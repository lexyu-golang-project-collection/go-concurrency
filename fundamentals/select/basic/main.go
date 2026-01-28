package main

import (
	"fmt"
	"time"
)

func selectDemo(done chan bool) {
	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(time.Second)
		ch1 <- "one"
	}()

	go func() {
		time.Sleep(time.Second)
		ch2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-ch1:
			fmt.Println("Received from ch1:", msg1)
		case msg2 := <-ch2:
			fmt.Println("Received from ch2:", msg2)
		case <-done:
			return
		}
	}
}

func main() {
	fmt.Println("\n3. Select Pattern:")
	done := make(chan bool)
	selectDemo(done)
}
