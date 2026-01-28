package main

import (
	"fmt"
	"sync"
)

func raceCondition() {
	fmt.Println("\nRace Condition Demo:")

	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++ // race condition
		}()
	}

	wg.Wait()
	fmt.Println("Counter value:", counter)
}

func main() {
	raceCondition()
}
