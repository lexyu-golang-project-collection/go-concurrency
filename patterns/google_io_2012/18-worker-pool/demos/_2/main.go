package main

import (
	"fmt"
	"math"
	"time"
)

// StockData 代表股票資料
type StockData struct {
	ID        int
	Symbol    string
	Price     float64
	Volume    int
	Timestamp time.Time
	Stages    []StageInfo
	Stage     string
	// 計算欄位
	MA5         float64 // 5日均價
	PriceChange float64 // 價格變動百分比
}

type StageInfo struct {
	Stage    string
	WorkerID int
	Time     time.Time
	Message  string // 加入處理訊息
}

type ETLPipeline struct {
	extractPool   *WorkerPool[StockData, StockData]
	transformPool *WorkerPool[StockData, StockData]
	loadPool      *WorkerPool[StockData, StockData]
}

func NewETLPipeline(workers, bufferSize int) *ETLPipeline {
	// Extract: 模擬從資料源讀取並清理資料
	extractWorker := func(workerID int, data StockData) StockData {
		time.Sleep(50 * time.Millisecond)

		// 模擬資料清理：移除極端值
		if data.Price > 1000 {
			data.Price = 1000
		}
		if data.Volume < 0 {
			data.Volume = 0
		}

		msg := fmt.Sprintf("Cleaned data - Price: %.2f, Volume: %d",
			data.Price, data.Volume)

		data.Stage = "extracted"
		data.Stages = append(data.Stages, StageInfo{
			Stage:    "extracted",
			WorkerID: workerID,
			Time:     time.Now(),
			Message:  msg,
		})

		fmt.Printf("Extract: ID=%d, Symbol=%s, WorkerID=%d, %s\n",
			data.ID, data.Symbol, workerID, msg)
		return data
	}

	// Transform: 計算技術指標
	transformWorker := func(workerID int, data StockData) StockData {
		time.Sleep(100 * time.Millisecond)

		// 模擬計算技術指標
		data.MA5 = data.Price * 0.95 // 模擬5日均價
		data.PriceChange = (data.Price - data.MA5) / data.MA5 * 100

		msg := fmt.Sprintf("Calculated MA5=%.2f, Change=%.2f%%",
			data.MA5, data.PriceChange)

		data.Stage = "transformed"
		data.Stages = append(data.Stages, StageInfo{
			Stage:    "transformed",
			WorkerID: workerID,
			Time:     time.Now(),
			Message:  msg,
		})

		fmt.Printf("Transform: ID=%d, Symbol=%s, WorkerID=%d, %s\n",
			data.ID, data.Symbol, workerID, msg)
		return data
	}

	// Load: 格式化資料準備儲存
	loadWorker := func(workerID int, data StockData) StockData {
		time.Sleep(30 * time.Millisecond)

		// 格式化數字，準備儲存
		data.Price = math.Round(data.Price*100) / 100
		data.MA5 = math.Round(data.MA5*100) / 100
		data.PriceChange = math.Round(data.PriceChange*100) / 100

		msg := fmt.Sprintf("Formatted for storage - Final Price=%.2f, MA5=%.2f, Change=%.2f%%",
			data.Price, data.MA5, data.PriceChange)

		data.Stage = "loaded"
		data.Stages = append(data.Stages, StageInfo{
			Stage:    "loaded",
			WorkerID: workerID,
			Time:     time.Now(),
			Message:  msg,
		})

		fmt.Printf("Load: ID=%d, Symbol=%s, WorkerID=%d, %s\n",
			data.ID, data.Symbol, workerID, msg)
		return data
	}

	return &ETLPipeline{
		extractPool:   NewWorkerPool[StockData, StockData](workers, bufferSize, extractWorker),
		transformPool: NewWorkerPool[StockData, StockData](workers, bufferSize, transformWorker),
		loadPool:      NewWorkerPool[StockData, StockData](workers, bufferSize, loadWorker),
	}
}

func (p *ETLPipeline) Start() {
	p.extractPool.Start()
	p.transformPool.Start()
	p.loadPool.Start()

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

func (p *ETLPipeline) Process(data StockData) {
	p.extractPool.Submit(data)
}

func (p *ETLPipeline) Results() <-chan StockData {
	return p.loadPool.Results()
}

func main() {
	pipeline := NewETLPipeline(3, 10)
	pipeline.Start()

	// 模擬股票資料
	symbols := []string{"AAPL", "GOOGL", "MSFT", "AMZN", "META"}

	go func() {
		for i := 0; i < 10; i++ {
			// 產生模擬資料
			data := StockData{
				ID:        i,
				Symbol:    symbols[i%len(symbols)],
				Price:     float64(100+i*10) + float64(i%3)*0.75,
				Volume:    1000 * (i + 1),
				Timestamp: time.Now(),
				Stages:    make([]StageInfo, 0),
			}
			pipeline.Process(data)
		}
		pipeline.extractPool.CloseQueue()
	}()

	processed := 0
	start := time.Now()

	fmt.Println("\n=== Processing Results ===")
	for result := range pipeline.Results() {
		processed++
		fmt.Printf("\nProcessed %s (ID=%d):\n", result.Symbol, result.ID)
		for _, stage := range result.Stages {
			fmt.Printf("  %s [Worker %d] @%s\n    %s\n",
				stage.Stage,
				stage.WorkerID,
				stage.Time.Format("15:04:05.000"),
				stage.Message)
		}
	}

	fmt.Printf("\nTotal time: %v\n", time.Since(start))
}
