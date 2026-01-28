package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type Job struct {
	ID       int
	Type     string
	Retry    int
	MaxRetry int
}

type Worker struct {
	ID       int
	JobChan  <-chan Job
	Result   chan<- string
	RetryOut chan<- Job
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("[Worker %d] Canceled\n", w.ID)
			return
		case job, ok := <-w.JobChan:
			if !ok {
				return
			}

			if rand.Float32() < 0.4 {
				if job.Retry < job.MaxRetry {
					job.Retry++
					fmt.Printf("Worker %d failed job %d, sending to retry\n", w.ID, job.ID)
					select {
					case w.RetryOut <- job:
					default:
						fmt.Printf("Retry buffer full, dropping job %d\n", job.ID)
					}
					continue
				} else {
					w.Result <- fmt.Sprintf("Worker %d failed job %d after max retries", w.ID, job.ID)
					continue
				}
			}

			time.Sleep(300 * time.Millisecond)
			w.Result <- fmt.Sprintf("Worker %d finished job %d", w.ID, job.ID)
		}
	}
}

func retryDispatcher(ctx context.Context, retryChan <-chan Job, jobQueue chan<- Job) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-retryChan:
			base := time.Duration(200*(1<<job.Retry)) * time.Millisecond
			jitter := time.Duration(rand.Intn(100)) * time.Millisecond
			time.Sleep(base + jitter)

			select {
			case jobQueue <- job:
				fmt.Printf("RetryDispatcher requeued job %d\n", job.ID)
			default:
				fmt.Printf("RetryDispatcher: job queue full, dropping job %d\n", job.ID)
			}
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	jobs := make(chan Job, 5)
	retryChan := make(chan Job, 5)
	results := make(chan string)

	for i := 1; i <= 3; i++ {
		w := &Worker{
			ID:       i,
			JobChan:  jobs,
			Result:   results,
			RetryOut: retryChan,
		}
		go w.Start(ctx)
	}

	go retryDispatcher(ctx, retryChan, jobs)

	go func() {
		jobID := 1
		for {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			default:
				baseDelay := 250 * time.Millisecond
				jitter := time.Duration(rand.Intn(100)) * time.Millisecond
				time.Sleep(baseDelay + jitter)

				job := Job{ID: jobID, Type: "TYPE-A", MaxRetry: 3}
				select {
				case jobs <- job:
					fmt.Println("Enqueued Job", jobID)
				default:
					fmt.Println("Job queue full, dropping Job", jobID)
				}
				jobID++
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("Result collector stopped:", ctx.Err())
			return
		case res := <-results:
			fmt.Println(res)
		}
	}
}
