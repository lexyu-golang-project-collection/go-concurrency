package main

import (
	"fmt"
	"sync"
	"time"
)

func waitGroupDemo() {
	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Worker %d doing work\n", id)
			time.Sleep(time.Second)
		}(i)
	}

	wg.Wait()
	fmt.Println("All workers done")
}

func main() {
	fmt.Println("\n4. WaitGroup Pattern:")
	waitGroupDemo()
}
