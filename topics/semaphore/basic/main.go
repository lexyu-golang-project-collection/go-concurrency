package main

import (
	"fmt"
	"sync"
	"time"
)

func semaphore() {
	fmt.Println("=== Semaphore Demo ===")

	concurrent := make(chan struct{}, 3) // limit: 3
	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			concurrent <- struct{}{} // acquire
			fmt.Printf("Worker %d is working\n", id)
			time.Sleep(time.Second)
			<-concurrent // release
		}(i)
	}

	wg.Wait()
}

func main() {
	semaphore()
}
