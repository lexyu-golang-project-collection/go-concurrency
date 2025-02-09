package main

import (
	"fmt"
)

func channelBasics() {
	// Unbuffered channel
	ch := make(chan string)
	go func() {
		ch <- "hello"
	}()
	msg := <-ch
	fmt.Println("Received:", msg)

	// Buffered channel
	bufCh := make(chan string, 2)
	bufCh <- "first"
	bufCh <- "second"
	fmt.Println(<-bufCh, <-bufCh)
}

func main() {
	fmt.Println("1. Channel Basics:")
	channelBasics()
}
