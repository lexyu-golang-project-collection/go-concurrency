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

func Google(query string) (result []Result) {
	result = append(result, Web(query))
	result = append(result, Image(query))
	result = append(result, Video(query))
	return result
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	start := time.Now()

	results := Google("golang")

	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
