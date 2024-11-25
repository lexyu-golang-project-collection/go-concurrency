package main

import (
	"fmt"
	"math/rand"
	"time"
)

func sender(message string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s =  %d", message, i)
			time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
		}
	}()

	return c
}

func main() {
	c := sender("Test")

	timeout := time.After(5 * time.Second)
	for {
		select {
		case str := <-c:
			fmt.Println(str)
		case <-timeout:
			fmt.Println("no response...")
			return
		}
	}
}
