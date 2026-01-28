package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string, ch chan string) {
	for i := 0; ; i++ {
		ch <- fmt.Sprintf("%s - %d", msg, i)
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func main() {

	ch := make(chan string)
	go boring("boring!", ch)

	for i := 1; i <= 7; i++ {
		fmt.Printf("You Say : %q\n", <-ch)
	}

	fmt.Println("Quit~~~~~~~~~~~~~")
}
