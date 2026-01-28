package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	symbol    string
	price     float64
	timestamp time.Time
	wait      chan bool
}

func fanIn(inputs ...<-chan Message) <-chan Message {
	c := make(chan Message)
	for i := range inputs {
		input := inputs[i]
		go func() {
			for {
				c <- <-input
			}

		}()
	}
	return c
}

func generateQuotes(symbol string, basePrice float64) <-chan Message {
	c := make(chan Message)
	waitForIt := make(chan bool)

	go func() {
		for {
			// Simulate price fluctuation
			priceChange := (rand.Float64() - 0.5) * 2
			newPrice := basePrice + priceChange

			c <- Message{
				symbol:    symbol,
				price:     newPrice,
				timestamp: time.Now(),
				wait:      waitForIt,
			}

			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
			<-waitForIt
		}
	}()
	return c
}

func main() {
	c := fanIn(
		generateQuotes("AAPL", 180.50),
		generateQuotes("GOOGL", 140.75),
	)
	for i := 0; i < 5; i++ {
		quote1 := <-c
		fmt.Printf("%s: $%.2f at %v\n",
			quote1.symbol, quote1.price, quote1.timestamp.Format("15:04:05"))

		quote2 := <-c
		fmt.Printf("%s: $%.2f at %v\n",
			quote2.symbol, quote2.price, quote2.timestamp.Format("15:04:05"))

		quote1.wait <- true
		quote2.wait <- true
	}

	fmt.Println("Market closed")
}
