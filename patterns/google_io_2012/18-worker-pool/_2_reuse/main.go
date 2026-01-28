package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// =============================================================
// Core Types and Interfaces
// =============================================================

type Job interface{}
type Result interface{}

// 包裝執行結果和錯誤
type JobResult[R Result] struct {
	Result R
	Error  error
}

// 函數現在可以返回錯誤
type Worker[T Job, R Result] func(ctx context.Context, workerID int, job T) (R, error)

// =============================================================
// Retry Policy Configuration
// =============================================================

// 定義重試策略
type RetryPolicy struct {
	MaxRetries      int           // 最大重試次數
	BaseDelay       time.Duration // 基礎延遲時間
	MaxDelay        time.Duration // 最大延遲時間
	BackoffStrategy BackoffStrategy
	JitterType      JitterType
}

type BackoffStrategy int

const (
	FixedBackoff BackoffStrategy = iota
	LinearBackoff
	ExponentialBackoff
)

type JitterType int

const (
	NoJitter JitterType = iota
	FullJitter
	EqualJitter
	DecorrelatedJitter
)

// =============================================================
// Worker Pool Structure and Internal Types
// =============================================================

type WorkerPool[T, R any] struct {
	// Configuration
	numOfWorkers int
	worker       Worker[T, R]
	retryPolicy  RetryPolicy

	// Channels
	jobs      chan T
	results   chan JobResult[R]
	retryJobs chan jobWithAttempt[T]

	// Concurrency Control
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	jobsOpen    atomic.Bool
	pendingJobs atomic.Int64 // 追蹤進行中的工作
}

type jobWithAttempt[T Job] struct {
	job     T
	attempt int
}

// =============================================================
// Constructor and Initialization
// =============================================================

func NewWorkerPool[T Job, R Result](
	numOfWorkers, bufferSize int,
	worker Worker[T, R],
	retryPolicy RetryPolicy,
) *WorkerPool[T, R] {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool[T, R]{
		// Configuration
		numOfWorkers: numOfWorkers,
		worker:       worker,
		retryPolicy:  retryPolicy,

		// Channels
		jobs:      make(chan T, bufferSize),
		results:   make(chan JobResult[R], bufferSize),
		retryJobs: make(chan jobWithAttempt[T], bufferSize),

		// Concurrency Control
		ctx:    ctx,
		cancel: cancel,
	}

	pool.jobsOpen.Store(true)
	return pool
}

// =============================================================
// Pool Lifecycle Management
// =============================================================

func (pool *WorkerPool[T, R]) Start() {
	// 啟動 workers
	for i := 1; i <= pool.numOfWorkers; i++ {
		pool.wg.Add(1)
		go pool.workerRoutine(i)
	}

	// 監控所有工作完成
	go pool.monitorCompletion()
}

func (pool *WorkerPool[T, R]) monitorCompletion() {
	// 等待所有 workers 完成
	pool.wg.Wait()
	close(pool.results)
}

func (pool *WorkerPool[T, R]) Stop() {
	pool.cancel()
}

func (pool *WorkerPool[T, R]) CloseQueue() {
	if pool.jobsOpen.CompareAndSwap(true, false) {
		close(pool.jobs)

		// 啟動超時機制，確保 workers 最終會退出
		go func() {
			// 等待一段時間讓 retry 完成
			timeout := time.Duration(pool.retryPolicy.MaxRetries) * pool.retryPolicy.MaxDelay * 2
			if timeout > 30*time.Second {
				timeout = 30 * time.Second
			}

			time.Sleep(timeout)

			// 如果還有待處理的工作，強制停止
			if pool.pendingJobs.Load() > 0 {
				pool.cancel()
			}
		}()
	}
}

// =============================================================
// Worker Execution Logic
// =============================================================

func (pool *WorkerPool[T, R]) workerRoutine(workerID int) {
	defer pool.wg.Done()

	for {
		select {
		case <-pool.ctx.Done():
			return
		case job, ok := <-pool.jobs:
			if !ok {
				// jobs channel 已關閉，但繼續處理 retry jobs
				pool.processRetryJobs(workerID)
				return
			}
			pool.processJob(workerID, job, 0)
		case jobAttempt, ok := <-pool.retryJobs:
			if !ok {
				return
			}
			pool.processJob(workerID, jobAttempt.job, jobAttempt.attempt)
		}
	}
}

func (pool *WorkerPool[T, R]) processRetryJobs(workerID int) {
	for {
		select {
		case <-pool.ctx.Done():
			return
		case jobAttempt, ok := <-pool.retryJobs:
			if !ok {
				return
			}
			pool.processJob(workerID, jobAttempt.job, jobAttempt.attempt)
		default:
			// 沒有更多 retry jobs，檢查是否還有其他工作
			if pool.pendingJobs.Load() == 0 {
				return
			}
			time.Sleep(10 * time.Millisecond) // 短暫等待
		}
	}
}

func (pool *WorkerPool[T, R]) processJob(workerID int, job T, attempt int) {
	pool.pendingJobs.Add(1)
	defer pool.pendingJobs.Add(-1)

	result, err := pool.worker(pool.ctx, workerID, job)

	if err != nil && attempt < pool.retryPolicy.MaxRetries {
		// 需要重試
		delay := pool.calculateDelay(attempt)

		// 非阻塞延遲
		select {
		case <-time.After(delay):
		case <-pool.ctx.Done():
			return
		}

		select {
		case pool.retryJobs <- jobWithAttempt[T]{job: job, attempt: attempt + 1}:
		case <-pool.ctx.Done():
			return
		}
		return
	}

	// 發送結果（成功或最終失敗）
	select {
	case pool.results <- JobResult[R]{Result: result, Error: err}:
	case <-pool.ctx.Done():
		return
	}
}

// =============================================================
// Retry and Backoff Algorithm Implementation
// =============================================================

func (pool *WorkerPool[T, R]) calculateDelay(attempt int) time.Duration {
	var delay time.Duration

	switch pool.retryPolicy.BackoffStrategy {
	case FixedBackoff:
		delay = pool.retryPolicy.BaseDelay
	case LinearBackoff:
		delay = time.Duration(attempt+1) * pool.retryPolicy.BaseDelay
	case ExponentialBackoff:
		delay = time.Duration(math.Pow(2, float64(attempt))) * pool.retryPolicy.BaseDelay
	}

	// 應用最大延遲限制
	if delay > pool.retryPolicy.MaxDelay {
		delay = pool.retryPolicy.MaxDelay
	}

	// 應用 Jitter
	return pool.applyJitter(delay, attempt)
}

func (pool *WorkerPool[T, R]) applyJitter(delay time.Duration, attempt int) time.Duration {
	switch pool.retryPolicy.JitterType {
	case NoJitter:
		return delay
	case FullJitter:
		// 隨機 0 到 delay 之間
		return time.Duration(rand.Int63n(int64(delay)))
	case EqualJitter:
		// 一半固定延遲 + 一半隨機延遲
		half := delay / 2
		return half + time.Duration(rand.Int63n(int64(half)))
	case DecorrelatedJitter:
		// AWS 建議的 decorrelated jitter
		min := int64(pool.retryPolicy.BaseDelay)
		max := int64(delay * 3)
		if max <= min {
			return pool.retryPolicy.BaseDelay
		}
		return time.Duration(min + rand.Int63n(max-min))
	default:
		return delay
	}
}

// =============================================================
// Public API Methods
// =============================================================

func (pool *WorkerPool[T, R]) Submit(job T) {
	if !pool.jobsOpen.Load() {
		return // 不接受新工作
	}

	select {
	case pool.jobs <- job:
	case <-pool.ctx.Done():
	}
}

func (pool *WorkerPool[T, R]) Results() <-chan JobResult[R] {
	return pool.results
}

// =============================================================
// Example Usage and Demo Types
// =============================================================

type DataJob struct {
	ID    int
	Value string
}

type DataResult struct {
	JobID     int
	Processed string
	WorkerID  int
}

func main() {
	// 建立重試策略
	retryPolicy := RetryPolicy{
		MaxRetries:      3,
		BaseDelay:       100 * time.Millisecond,
		MaxDelay:        5 * time.Second,
		BackoffStrategy: ExponentialBackoff,
		JitterType:      EqualJitter,
	}

	// 建立會偶爾失敗的工作處理函數
	worker := func(ctx context.Context, workerID int, job DataJob) (DataResult, error) {
		// 模擬 20% 的失敗率
		if rand.Float32() < 0.2 {
			return DataResult{}, fmt.Errorf("simulated failure for job %d", job.ID)
		}

		// 模擬處理時間
		time.Sleep(50 * time.Millisecond)

		return DataResult{
			JobID:     job.ID,
			Processed: fmt.Sprintf("Processed-%s", job.Value),
			WorkerID:  workerID,
		}, nil
	}

	// 建立工作池
	pool := NewWorkerPool[DataJob, DataResult](5, 100, worker, retryPolicy)
	pool.Start()

	// 提交工作
	go func() {
		for i := 0; i < 50; i++ {
			pool.Submit(DataJob{
				ID:    i,
				Value: fmt.Sprintf("Data-%d", i),
			})
		}
		pool.CloseQueue()
	}()

	// 處理結果
	successCount := 0
	failureCount := 0

	for result := range pool.Results() {
		if result.Error != nil {
			fmt.Printf("❌ Final failure: JobID=%d, Error=%s\n",
				result.Result.JobID, result.Error)
			failureCount++
		} else {
			fmt.Printf("✅ Success: JobID=%d, Processed=%s, WorkerID=%d\n",
				result.Result.JobID, result.Result.Processed, result.Result.WorkerID)
			successCount++
		}
	}

	fmt.Printf("\n📊 Results: %d successful, %d failed\n", successCount, failureCount)
}
