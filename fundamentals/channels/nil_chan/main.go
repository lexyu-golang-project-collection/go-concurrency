package main

import (
	"fmt"
	"time"
)

func nilChannelOperations() {
	fmt.Println("=== Nil Channel Forever Block ===")

	var nilCh chan int // nil channel

	// 在 goroutine 中操作 nil channel
	go func() {
		fmt.Println("Trying to read from nil channel (will block forever)")
		<-nilCh
	}()

	// 主程式持續運行
	for {
		fmt.Println("Main program is still running...")
		time.Sleep(time.Second)
	}
}

func main() {
	nilChannelOperations()
}
