package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	numberOfJobs    = 12
	numberOfWorkers = 3
)

func workerUnefficient(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "fnished job", j)
		results <- j * 2
	}
}

func workers(workerId int, jobs <-chan int, result chan<- int) {
	var wg sync.WaitGroup

	for job := range jobs {
		wg.Add(1)
		go func(job int) {
			fmt.Println("worker ", workerId, "started job ", job)
			time.Sleep(time.Second)
			fmt.Println("worker ", workerId, "finished job ", job)
			result <- job * 3
			wg.Done()
		}(job)
	}
	wg.Wait()
}

func main() {
	// Demo1 和 Demo2 的對比測試
	// compareDemo()

	// 或者您想單獨執行 Demo
	runDemo1()
	// runDemo2()
}

// 比較兩種實現方式的效能差異
func compareDemo() {
	fmt.Println("====== 開始效能比較測試 ======")
	fmt.Printf("Initial goroutines: %d\n\n", runtime.NumGoroutine())

	// 執行 Demo1
	fmt.Println("--- 執行 Demo1 (每個 worker 一個 goroutine) ---")
	startMonitoring("Demo1")
	start := time.Now()
	Demo1()
	fmt.Printf("Demo1 執行時間: %v\n", time.Since(start))
	fmt.Printf("Demo1 結束後的 goroutines 數量: %d\n\n", runtime.NumGoroutine())

	time.Sleep(time.Second) // 確保前一個 demo 完全結束

	// 執行 Demo2
	fmt.Println("--- 執行 Demo2 (每個任務一個 goroutine) ---")
	startMonitoring("Demo2")
	start = time.Now()
	Demo2()
	fmt.Printf("Demo2 執行時間: %v\n", time.Since(start))
	fmt.Printf("Demo2 結束後的 goroutines 數量: %d\n", runtime.NumGoroutine())

	fmt.Println("====== 比較測試結束 ======")
}

// 監控 goroutine 數量的函數
func startMonitoring(demoName string) {
	go func() {
		for i := 0; i < 5; i++ {
			fmt.Printf("Current goroutines during %s: %d\n",
				demoName, runtime.NumGoroutine())
			time.Sleep(time.Millisecond * 200)
		}
	}()
	// 給監控 goroutine 一點時間先啟動
	time.Sleep(time.Millisecond * 50)
}

// 單獨運行 Demo1
func runDemo1() {
	fmt.Println("====== 執行 Demo1 ======")
	fmt.Printf("Initial goroutines: %d\n", runtime.NumGoroutine())
	startMonitoring("Demo1")
	Demo1()
	fmt.Printf("Final goroutines: %d\n", runtime.NumGoroutine())
}

// 單獨運行 Demo2
func runDemo2() {
	fmt.Println("====== 執行 Demo2 ======")
	fmt.Printf("Initial goroutines: %d\n", runtime.NumGoroutine())
	startMonitoring("Demo2")
	Demo2()
	fmt.Printf("Final goroutines: %d\n", runtime.NumGoroutine())
}

// Demo1: 每個 worker 一個 goroutine
func Demo1() {
	jobs := make(chan int, numberOfJobs)
	results := make(chan int, numberOfJobs)

	// 啟動 workers
	for w := 1; w <= numberOfWorkers; w++ {
		fmt.Printf("啟動 worker %d\n", w)
		go workerUnefficient(w, jobs, results)
	}

	// 發送工作
	for j := 1; j <= numberOfJobs; j++ {
		fmt.Printf("發送工作 %d 到 channel\n", j)
		jobs <- j
	}
	close(jobs)

	// 收集結果
	for a := 1; a <= numberOfJobs; a++ {
		fmt.Printf("收到結果: %d\n", <-results)
	}
	close(results)
}

// Demo2: 每個任務一個 goroutine
func Demo2() {
	jobs := make(chan int, numberOfJobs)
	results := make(chan int, numberOfJobs)

	// 啟動 workers
	for w := 1; w <= numberOfWorkers; w++ {
		fmt.Printf("啟動 worker %d\n", w)
		go workers(w, jobs, results)
	}

	// 發送工作
	for j := 1; j <= numberOfJobs; j++ {
		fmt.Printf("發送工作 %d 到 channel\n", j)
		jobs <- j
	}
	close(jobs)

	// 收集結果
	for a := 1; a <= numberOfJobs; a++ {
		fmt.Printf("收到結果: %d\n", <-results)
	}
	close(results)
}
