package main

import (
	"fmt"
	"net/http"
	"time"

	"../log"
)

func main() {

	http.HandleFunc("/", log.Decorate(handler))
	panic(http.ListenAndServe("127.0.0.1:8080", nil))

}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// ctx = context.WithValue(ctx, int(42), int64(100))
	log.Println(ctx, "handler started")
	defer log.Println(ctx, "handler ended")

	fmt.Printf("value foo for is %v\n", ctx.Value("foo"))

	select {
	case <-time.After(5 * time.Second):
		// time.Sleep(5 * time.Second)
		fmt.Fprintln(w, "hello there!")
	case <-ctx.Done():
		err := ctx.Err()
		log.Println(ctx, err.Error())
		http.Error(w, ctx.Err().Error(), http.StatusInternalServerError)
	}
}
