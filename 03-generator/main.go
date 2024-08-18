package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) <-chan string {
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

	J := boring("Jack")
	K := boring("K")

	for i := 1; i <= 10; i++ {
		fmt.Println(<-J)
		fmt.Println(<-K)
	}

	// or use for range
	// for msg := range K {
	// 	fmt.Println(msg)
	// }

	// for J != nil || K != nil {
	// 	select {
	// 	case msg, ok := <-J:
	// 		if ok {
	// 			fmt.Println(msg)
	// 		} else {
	// 			J = nil
	// 		}
	// 	case msg, ok := <-K:
	// 		if ok {
	// 			fmt.Println(msg)
	// 		} else {
	// 			K = nil
	// 		}
	// 	}
	// }

	fmt.Println("Quit~~~~~~~~~~~~~")
}
