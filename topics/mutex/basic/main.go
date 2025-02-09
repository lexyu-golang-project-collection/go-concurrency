package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeCounter struct {
	mu    sync.Mutex
	count int
}

func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func main() {
	fmt.Println("\n5. Mutex Demo:")
	counter := SafeCounter{}
	for i := 0; i < 3; i++ {
		go counter.Increment()
	}
	time.Sleep(time.Second)
	fmt.Println("Counter:", counter.Value())
}
