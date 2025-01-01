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
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}

}

func Google(query string) []Result {
	resultsChan := make(chan Result)

	go func() {
		resultsChan <- Web(query)
	}()
	go func() {
		resultsChan <- Image(query)
	}()
	go func() {
		resultsChan <- Video(query)
	}()

	var results []Result

	// with timeout and select
	timeout := time.After(50 * time.Millisecond)

	for i := 0; i < 3; i++ {
		select {
		case r := <-resultsChan:
			results = append(results, r)
		case <-timeout:
			fmt.Println("timeout!")
			return results
		}
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
