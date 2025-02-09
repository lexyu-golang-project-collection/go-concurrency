package main

import (
	"fmt"
	"sync"
	"time"
)

func waitGroupMisuse1() {
	fmt.Println("\nWaitGroup Misuse Negative counter Demo:")

	var wg sync.WaitGroup

	// Case 1: Negative counter
	fmt.Println("Case 1: Calling Done without Add")
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered from:", r)
			}
		}()
		wg.Done() // panic: negative WaitGroup counter
	}()

	time.Sleep(time.Second)
}

func waitGroupMisuse2() {
	fmt.Println("\nWaitGroup Misuse Missing Done Demo:")

	var wg sync.WaitGroup

	fmt.Println("Case 2: Missing Done call")
	wg.Add(1)
	go func() {
		fmt.Println("Goroutine working...")
		// wg.Done() is missing
	}()

	time.Sleep(time.Second)
}

func waitGroupMisuse3() {
	fmt.Println("\nWaitGroup Misuse Wrong Add count Demo:")

	var wg sync.WaitGroup
	fmt.Println("Case 3: Incorrect Add count")
	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Println("Only one goroutine working...")
	}()

	time.Sleep(time.Second)
}

func main() {
	// waitGroupMisuse1()
	waitGroupMisuse2()
	// waitGroupMisuse3()
}
