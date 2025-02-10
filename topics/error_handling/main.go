package main

import (
	"fmt"
	"time"
)

type Result struct {
	Value string
	Err   error
}

func errHandle() {
	fmt.Println("\n=== Error Handling Pattern ===")

	ch := make(chan Result)
	go func() {
		// Simulate some work that might error
		time.Sleep(time.Second)
		if time.Now().UnixNano()%2 == 0 {
			ch <- Result{Err: fmt.Errorf("random error")}
		} else {
			ch <- Result{Value: "success"}
		}
	}()

	// Handle result
	if result := <-ch; result.Err != nil {
		fmt.Printf("Error occurred: %v\n", result.Err)
	} else {
		fmt.Printf("Success: %s\n", result.Value)
	}
}

func main() {
	errHandle()
}
