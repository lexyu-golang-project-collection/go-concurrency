package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

func memoryLeak1() {
	fmt.Println("\n=== Memory Leak ===")

	// Goroutine leak
	for i := 0; i < 10; i++ {
		go func() {
			// This goroutine will never exit
			select {}
		}()
	}

	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())

	time.Sleep(time.Second)
}

func memoryLeak2() {
	fmt.Println("\n=== Memory Leak ===")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Channel leak
	for i := 0; i < 10; i++ {
		ch := make(chan int)
		go func() {
			// Channel is never closed or read from
			<-ch
		}()
	}
	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())

	<-sigChan
}

func memoryLeak3() {
	fmt.Println("\n=== Memory Leak ===")

	// Timer leak
	for i := 0; i < 10; i++ {
		timer := time.NewTimer(time.Hour)
		fmt.Println("timer ", timer)
		timer.Stop() // Forgot to stop timer
	}

	fmt.Printf("Number of goroutines: %d\n", runtime.NumGoroutine())

	time.Sleep(time.Second)
}

func main() {
	//memoryLeak1()
	//memoryLeak2()
	memoryLeak3()
}
