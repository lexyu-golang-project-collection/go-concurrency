package main

import (
	"context"
	"time"
)

func main() {
	ctx := context.Background()
	hub := NewHub()

	sub1 := NewSubscriber("sub-01")
	sub2 := NewSubscriber("sub-02")
	sub3 := NewSubscriber("sub-03")

	hub.subscribe(ctx, sub1)
	hub.subscribe(ctx, sub2)
	hub.subscribe(ctx, sub3)

	_ = hub.publish(ctx, &Message{data: []byte("test-01")})
	_ = hub.publish(ctx, &Message{data: []byte("test-02")})
	_ = hub.publish(ctx, &Message{data: []byte("test-03")})
	time.Sleep(2 * time.Second)

	hub.unSubscribe(ctx, sub3)

	_ = hub.publish(ctx, &Message{data: []byte("test-04")})
	_ = hub.publish(ctx, &Message{data: []byte("test-05")})

	time.Sleep(2 * time.Second)
}
