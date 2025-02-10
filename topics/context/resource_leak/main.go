package main

import (
	"context"
	"fmt"
	"time"
)

func contextErrors() {
	fmt.Println("\n=== Context Errors ===")

	// Using cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Start work after context is already cancelled
	go func() {
		select {
		case <-ctx.Done():
			fmt.Println("Context was already cancelled")
		default:
			fmt.Println("This won't be printed")
		}
	}()

	// Forgot to cancel context
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second)
	//defer cancel2() // Forgot to cancel - potential resource leak

	go func() {
		<-ctx2.Done()
		fmt.Println("Context timed out but wasn't properly cancelled")
	}()

	time.Sleep(2 * time.Second)
	cancel2() // Late resource_leak
}

func main() {
	contextErrors()
}
