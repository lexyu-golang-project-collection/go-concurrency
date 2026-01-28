package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Ball struct {
	hits int
}

func player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++
		fmt.Println(name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}
}

func main() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	n := r.Int63n(10)
	table := make(chan *Ball)
	go player("A player ping", table)
	go player("B player pong", table)
	table <- &Ball{hits: 0}
	time.Sleep(time.Duration(n) * time.Second)
	<-table
	panic("stack")
}
