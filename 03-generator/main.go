package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generator(msg string) <-chan string {
	ch := make(chan string)
	go func() {
		for i := 1; i <= 10; i++ {
			ch <- fmt.Sprintf("%s - %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
		// The sender should close the channel
		close(ch)
	}()

	return ch
}

func main() {

	J := generator("Jack")
	K := generator("K")

	// 1.
	// Method1(J, K)

	// 2.
	Method2(K, J)

	// 3.
	// Method3(J, K)

	fmt.Println("Quit~~~~~~~~~~~~~")
}

// for loop
func Method1(ch1, ch2 <-chan string) {
	for i := 1; i <= 10; i++ {
		fmt.Println(<-ch1)
		fmt.Println(<-ch2)
	}
}

// for range
func Method2(channels ...<-chan string) {
	for _, ch := range channels {
		for msg := range ch {
			fmt.Println(msg)
		}
	}
}

// select 和 nil 通道
func Method3(ch1, ch2 <-chan string) {
	for ch1 != nil || ch2 != nil {
		select {
		case msg, ok := <-ch1:
			if ok {
				fmt.Println(msg)
			} else {
				ch1 = nil
			}
		case msg, ok := <-ch2:
			if ok {
				fmt.Println(msg)
			} else {
				ch2 = nil
			}
		}
	}
}
