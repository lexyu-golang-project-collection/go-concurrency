package main

import (
	"context"
	"log"
	"time"
)

func sleepAndTalk(ctx context.Context, duration time.Duration, msg string) {
	select {
	case <-time.After(duration):
		log.Println(msg)
	case <-ctx.Done():
		log.Panicln(ctx.Err())
	}
}

func main() {
	log.Println("started")

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	time.AfterFunc(time.Second, cancel)
	sleepAndTalk(ctx, 5*time.Second, "hello")
}
