package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("viedo")
)

type Result string
type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}

}

func Google(query string) []Result {
	resChan := make(chan Result)

	go func() {
		resChan <- Web(query)
	}()
	go func() {
		resChan <- Image(query)
	}()
	go func() {
		resChan <- Video(query)
	}()

	var results []Result
	for i := 0; i < 3; i++ {
		results = append(results, <-resChan)
	}
	return results
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	start := time.Now()
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
