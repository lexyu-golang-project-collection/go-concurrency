package main

import (
	"context"
	"fmt"
	"time"
)

func executeWithTimeout(ctx context.Context) error {
	fmt.Println("Starting tasks with 2 second timeout...")

	timeoutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer func() {
		cancel()
		fmt.Println("Context cancelled")
	}()

	errCh := make(chan error, 2)

	fmt.Println("Launching subTask1...")
	go func() {
		if err := subTask1(timeoutCtx); err != nil {
			fmt.Printf("subTask1 error: %v\n", err)
			errCh <- fmt.Errorf("subTask1 failed: %v", err)
			return
		}
		fmt.Println("subTask1 completed successfully")
		errCh <- nil
	}()

	fmt.Println("Launching subTask2...")
	go func() {
		if err := subTask2(timeoutCtx); err != nil {
			fmt.Printf("subTask2 error: %v\n", err)
			errCh <- fmt.Errorf("subTask2 failed: %v", err)
			return
		}
		fmt.Println("subTask2 completed successfully")
		errCh <- nil
	}()

	fmt.Println("Waiting for tasks to complete...")
	for i := 0; i < 2; i++ {
		if err := <-errCh; err != nil {
			fmt.Printf("Task error received: %v\n", err)
			return err
		}
	}

	fmt.Println("All tasks completed")
	return nil
}

func subTask1(ctx context.Context) error {
	fmt.Println("subTask1: Starting 1 second operation")
	select {
	case <-time.After(1 * time.Second):
		fmt.Println("subTask1: Completed normally")
		return nil
	case <-ctx.Done():
		fmt.Printf("subTask1: Interrupted - %v\n", ctx.Err())
		return ctx.Err()
	}
}

func subTask2(ctx context.Context) error {
	fmt.Println("subTask2: Starting 3 second operation")
	select {
	case <-time.After(3 * time.Second):
		fmt.Println("subTask2: Completed normally")
		return nil
	case <-ctx.Done():
		fmt.Printf("subTask2: Interrupted - %v\n", ctx.Err())
		return ctx.Err()
	}
}

func main() {
	ctx := context.Background()
	if err := executeWithTimeout(ctx); err != nil {
		fmt.Printf("Main: Execution failed - %v\n", err)
	}
}
