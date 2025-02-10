package main

import (
	"fmt"
	"time"
)

func selectErrors() {
	fmt.Println("=== Select Forever Block ===")

	// 使用 goroutine 執行 empty select
	go func() {
		fmt.Println("Empty select will block forever")
		select {}
	}()

	for {
		fmt.Println("Main program is still running...")
		time.Sleep(time.Second)
	}
}

func main() {
	selectErrors()
}
