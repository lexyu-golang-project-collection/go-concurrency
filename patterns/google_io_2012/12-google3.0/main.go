package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web1   = fakeSearch("web1")
	Web2   = fakeSearch("web2")
	Image1 = fakeSearch("image1")
	Image2 = fakeSearch("image2")
	Video1 = fakeSearch("viedo1")
	Video2 = fakeSearch("viedo2")
)

type Result string
type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}

}

// avoid discarding result from slow server
func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	for i := range replicas {
		go func(idx int) {
			c <- replicas[idx](query)
		}(i)
	}
	return <-c
}

func Google(query string) []Result {
	resultsChan := make(chan Result)

	go func() {
		resultsChan <- First(query, Web1, Web2)
	}()
	go func() {
		resultsChan <- First(query, Image1, Image2)
	}()
	go func() {
		resultsChan <- First(query, Video1, Video2)
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
