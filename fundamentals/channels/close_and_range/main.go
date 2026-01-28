package main

import "fmt"

func channelCloseRange() {
	ch := make(chan int, 3)

	// Producer
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
		}
		close(ch)
	}()

	// Consumer
	for val := range ch {
		fmt.Println("Value:", val)
	}

	// Reading from closed channel
	val, ok := <-ch
	fmt.Printf("Read closed channel: val=%v, ok=%v\n", val, ok)
}

func main() {
	fmt.Println("\n2. Channel Close & Range:")
	channelCloseRange()
}
