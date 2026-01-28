package main

import (
	"fmt"
	"sync"
	"time"
)

func fanOut() {
	work := make(chan int)
	done := make(chan bool)

	// Producer
	go func() {
		for i := 1; i <= 10; i++ {
			work <- i
		}
		close(work)
	}()

	// Start 3 workers
	var wg sync.WaitGroup
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for w := range work {
				fmt.Printf("Worker %d processing %d\n", id, w)
				time.Sleep(time.Millisecond * 100)
			}
		}(i)
	}

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		done <- true
	}()

	<-done
}

func main() {
	fmt.Println("\n8. Fan-out Pattern:")
	fanOut()
}
