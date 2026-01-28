package main

import (
	"context"
	"fmt"
	"time"
)

func contextDemo() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	go worker(ctx, "worker1")
	go worker(ctx, "worker2")

	<-ctx.Done()
	fmt.Println("Main: context done")
}

func worker(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s: stopping\n", name)
			return
		default:
			fmt.Printf("%s: working\n", name)
			time.Sleep(time.Second)
		}
	}
}

func main() {
	fmt.Println("\n6. Context Pattern:")
	contextDemo()
}
