package main

import (
	"fmt"
	"math/rand"
	"time"
)

func baker(name string, quit chan string, r *rand.Rand) <-chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		for i := 1; ; i++ {
			select {
			case c <- fmt.Sprintf("%s: Cake #%d is ready!", name, i):
			case message := <-quit:
				fmt.Printf("[%s] Stopping work: %s. Cleaning up...\n", name, message)
				quit <- fmt.Sprintf("%s: All baking equipment is cleaned up. Goodbye!", name)
				return
			}
			time.Sleep(time.Duration(r.Intn(800)+200) * time.Millisecond)
		}
	}()
	return c
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	quit := make(chan string)
	bakery := baker("Chef Paul", quit, r)

	for i := 0; i < 5; i++ {
		fmt.Println(<-bakery)
	}

	quit <- "Stop baking"
	fmt.Println("Final Message:", <-quit)
}
