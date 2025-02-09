package main

import (
	"fmt"
	"sync"
	"time"
)

func mutexDeadlocks() {
	fmt.Println("\nMutex Deadlock Demo:")

	var mutex1, mutex2 sync.Mutex
	var wg sync.WaitGroup

	// Goroutine 1: locks mutex1 then mutex2
	wg.Add(1)
	go func() {
		defer wg.Done()
		mutex1.Lock()
		fmt.Println("Goroutine 1: locked mutex1")
		time.Sleep(100 * time.Millisecond)

		mutex2.Lock()
		fmt.Println("Goroutine 1: locked mutex2")

		mutex2.Unlock()
		mutex1.Unlock()
	}()

	// Goroutine 2: locks mutex2 then mutex1
	wg.Add(1)
	go func() {
		defer wg.Done()
		mutex2.Lock()
		fmt.Println("Goroutine 2: locked mutex2")
		time.Sleep(100 * time.Millisecond)

		mutex1.Lock()
		fmt.Println("Goroutine 2: locked mutex1")

		mutex1.Unlock()
		mutex2.Unlock()
	}()

	wg.Wait()
}

func main() {
	mutexDeadlocks()
}
