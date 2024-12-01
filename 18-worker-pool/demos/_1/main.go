package main

import (
	"fmt"
	"time"
)

// ETLData 代表要處理的資料
type ETLData struct {
	ID     int
	Data   string
	Stages []StageInfo
	Stage  string
}

// StageInfo 記錄每個處理階段的資訊
type StageInfo struct {
	Stage    string
	WorkerID int
	Time     time.Time
}

// ETLPipeline 代表 ETL 處理管線
type ETLPipeline struct {
	extractPool   *WorkerPool[ETLData, ETLData]
	transformPool *WorkerPool[ETLData, ETLData]
	loadPool      *WorkerPool[ETLData, ETLData]
}

// NewETLPipeline 建立新的 ETL 管線
func NewETLPipeline(workers, bufferSize int) *ETLPipeline {
	// Extract worker
	extractWorker := func(workerID int, data ETLData) ETLData {
		time.Sleep(100 * time.Millisecond)
		data.Stage = "extracted"
		data.Stages = append(data.Stages, StageInfo{
			Stage:    "extracted",
			WorkerID: workerID,
			Time:     time.Now(),
		})
		fmt.Printf("Extract: ID=%d, WorkerID=%d\n", data.ID, workerID)
		return data
	}

	// Transform worker
	transformWorker := func(workerID int, data ETLData) ETLData {
		time.Sleep(150 * time.Millisecond)
		data.Stage = "transformed"
		data.Stages = append(data.Stages, StageInfo{
			Stage:    "transformed",
			WorkerID: workerID,
			Time:     time.Now(),
		})
		fmt.Printf("Transform: ID=%d, WorkerID=%d\n", data.ID, workerID)
		return data
	}

	// Load worker
	loadWorker := func(workerID int, data ETLData) ETLData {
		time.Sleep(80 * time.Millisecond)
		data.Stage = "loaded"
		data.Stages = append(data.Stages, StageInfo{
			Stage:    "loaded",
			WorkerID: workerID,
			Time:     time.Now(),
		})
		fmt.Printf("Load: ID=%d, WorkerID=%d\n", data.ID, workerID)
		return data
	}

	return &ETLPipeline{
		extractPool:   NewWorkerPool[ETLData, ETLData](workers, bufferSize, extractWorker),
		transformPool: NewWorkerPool[ETLData, ETLData](workers, bufferSize, transformWorker),
		loadPool:      NewWorkerPool[ETLData, ETLData](workers, bufferSize, loadWorker),
	}
}

// Start 啟動 ETL 管線
func (p *ETLPipeline) Start() {
	p.extractPool.Start()
	p.transformPool.Start()
	p.loadPool.Start()

	// 連接各階段的處理管線
	go func() {
		for data := range p.extractPool.Results() {
			p.transformPool.Submit(data)
		}
		p.transformPool.CloseQueue()
	}()

	go func() {
		for data := range p.transformPool.Results() {
			p.loadPool.Submit(data)
		}
		p.loadPool.CloseQueue()
	}()
}

// Process 處理資料
func (p *ETLPipeline) Process(data ETLData) {
	p.extractPool.Submit(data)
}

// Results 取得最終結果
func (p *ETLPipeline) Results() <-chan ETLData {
	return p.loadPool.Results()
}

func main() {
	// 建立 ETL 管線
	pipeline := NewETLPipeline(5, 10)
	pipeline.Start()

	// 提交資料
	go func() {
		for i := 1; i <= 100; i++ {
			pipeline.Process(ETLData{
				ID:     i,
				Data:   fmt.Sprintf("RawData-%d", i),
				Stages: make([]StageInfo, 0),
			})
		}
		pipeline.extractPool.CloseQueue()
	}()

	// 處理結果
	processed := 0
	start := time.Now()

	for result := range pipeline.Results() {
		processed++
		fmt.Printf("\nComplete Processing ID=%d:\n", result.ID)
		for _, stage := range result.Stages {
			fmt.Printf("  - Stage: %s, WorkerID: %d, Time: %s\n",
				stage.Stage, stage.WorkerID, stage.Time.Format("15:04:05.000"))
		}
	}

	fmt.Printf("\nTotal time: %v\n", time.Since(start))
}
