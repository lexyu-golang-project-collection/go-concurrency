package main

import (
	"fmt"
	"sync/atomic"
)

// 用於追蹤 goroutine 編號
var counter uint64

func f(id uint64, left, right chan int) {
	rightValue := <-right
	newValue := 1 + rightValue
	fmt.Printf("Goroutine %d: received %d from right, sending %d to left\n",
		id, rightValue, newValue)
	left <- newValue
}

func main() {
	const n = 10
	leftmost := make(chan int)
	left := leftmost
	right := leftmost

	fmt.Println("Creating chain of", n, "goroutines...")

	for i := 0; i < n; i++ {
		right = make(chan int)
		id := atomic.AddUint64(&counter, 1)
		fmt.Printf("Creating goroutine %d\n", id)
		go f(id, left, right)
		left = right
	}

	fmt.Println("\nStarting chain reaction by sending 1 to rightmost channel...")
	go func(c chan int) {
		fmt.Println("Anonymous goroutine: sending 1 to rightmost")
		c <- 1
	}(right)

	result := <-leftmost
	fmt.Printf("\nFinal value received at leftmost: %d\n", result)
}
