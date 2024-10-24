package main

import (
	"fmt"
	"math/rand"
	"time"
)

func producer(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s = %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()

	return c
}

// Merge Channels
func fanIn(ch1, ch2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			v1 := <-ch1
			c <- v1
		}
	}()
	go func() {
		for {
			c <- <-ch2
		}
	}()
	return c
}

func fanInSimple(cs ...<-chan string) <-chan string {
	c := make(chan string)

	for k, v := range cs {
		fmt.Println("k=", k)
		go func(cv <-chan string) {
			for {
				c <- <-cv
			}
		}(v)
	}

	return c
}

func main() {
	// Demo1()

	Demo2()
}

func Demo1() {
	c := fanIn(producer("test-channel-1"), producer("test-channel-2"))
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
}

func Demo2() {
	c2 := fanInSimple(producer("test-channel-1"), producer("test-channel-2"))
	for i := 0; i < 10; i++ {
		fmt.Println(<-c2)
	}
}
