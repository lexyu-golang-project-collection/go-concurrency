package main

import (
	"fmt"
	"time"
)

func producer(ch chan<- int, count int) {
	for i := 1; i <= count; i++ {
		ch <- i
	}
	close(ch)
}

func consumer(ch <-chan int, processingTime time.Duration) <-chan bool {
	done := make(chan bool)
	go func() {
		for job := range ch {
			fmt.Printf("Processing job %d\n", job)
			time.Sleep(processingTime)
		}
		done <- true
	}()
	return done
}

func pubsub(jobCount int, bufferSize int, processingTime time.Duration) {
	start := time.Now()

	jobs := make(chan int, bufferSize)

	go producer(jobs, jobCount)

	done := consumer(jobs, processingTime)
	<-done

	fmt.Printf("Time taken: %v\n", time.Since(start))
}

func main() {
	pubsub(
		500,                 // 工作數量
		100,                 // buffer 大小
		10*time.Millisecond, // 處理時間
	)
}
