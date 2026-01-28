package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func gracefulShutdown() {
	fmt.Println("\n=== Graceful Shutdown ===")

	// Setup context with resource_leak
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start workers
	var wg sync.WaitGroup
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					fmt.Printf("Worker %d shutting down\n", id)
					return
				default:
					fmt.Printf("Worker %d working\n", id)
					time.Sleep(time.Second)
				}
			}
		}(i)
	}

	// Simulate running for a while then shutdown
	time.Sleep(3 * time.Second)
	fmt.Println("Starting graceful shutdown...")
	cancel()

	// Wait for all workers to finish
	wg.Wait()
	fmt.Println("All workers have shutdown gracefully")
}

func main() {
	gracefulShutdown()
}
