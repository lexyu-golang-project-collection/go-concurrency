package main

import (
	"fmt"
	"runtime/debug"
	"time"
)

// Demo 1: Send to unbuffered channel without receiver
func demo1() {
	fmt.Println("=== Demo 1: Send without receiver ===")

	ch := make(chan string)

	ch <- "message" //  deadlock！

	fmt.Println("Message sent")
}

// Demo 2: Receive from unbuffered channel without sender
func demo2() {
	fmt.Println("=== Demo 2: Receive without sender ===")

	ch := make(chan string)

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Demo panicked: %v\n", r)
			fmt.Printf("Stack trace:\n%s\n", debug.Stack())
		}
	}()
	msg := <-ch // deadlock!
	fmt.Println(msg)

	time.Sleep(time.Second)
}

// Demo 3: Channel buffer overflow
func demo3() {
	fmt.Println("=== Demo 3: Buffer overflow ===")
	ch := make(chan string, 1)
	ch <- "first" // OK
	fmt.Println("First message sent")
	ch <- "second" // deadlock
}

// Demo 4: Operations on closed channel
func demo4() {
	fmt.Println("=== Demo 4: Operations on closed channel ===")
	ch := make(chan int)
	close(ch)

	// 從已關閉的 channel 讀取 (不會 deadlock，會得到零值)
	val := <-ch
	fmt.Printf("Read from closed channel: %v\n", val)

	// 寫入已關閉的 channel (會 panic)
	ch <- 1 // panic: send on closed channel
}

func main() {
	// demo1()
	// demo2()
	// demo3()
	demo4()
}
