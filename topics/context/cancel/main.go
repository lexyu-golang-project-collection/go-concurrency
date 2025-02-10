package main

import (
	"context"
	"fmt"
	"time"
)

func cancelExample() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go worker(ctx, "worker-1")
	go worker(ctx, "worker-2")

	// 模擬某些條件下需要取消
	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(time.Second) // 等待 worker 收到取消信號
}

func worker(ctx context.Context, id string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %s: Cancelled\n", id)
			return
		case <-time.After(500 * time.Millisecond):
			fmt.Printf("Worker %s: Working...\n", id)
		}
	}
}

func main() {
	cancelExample()
}
