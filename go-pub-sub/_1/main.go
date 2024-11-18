package main

import (
	"context"
	"fmt"
	"time"
)

func main() {

	ctx := context.Background()
	hub1 := NewHub()
	hub2 := NewHub()
	publisher := NewPublisher(hub1, hub2)

	sub1 := NewSubscriber("sub-01")
	hub1.subscribe(ctx, sub1)

	sub2 := NewSubscriber("sub-02")
	hub2.subscribe(ctx, sub2)

	sub3 := NewSubscriber("sub-03")
	hub2.subscribe(ctx, sub3)

	fmt.Println("Hub1 Users = ", hub1.getUsers())
	fmt.Println("Hub2 Users = ", hub2.getUsers())

	publisher.publish(ctx, &Message{data: []byte("test-01")})
	publisher.publish(ctx, &Message{data: []byte("test-02")})
	publisher.publish(ctx, &Message{data: []byte("test-03")})

	time.Sleep(2 * time.Second)

	hub2.unSubscribe(sub3)

	fmt.Println("Hub2 Users = ", hub2.getUsers())

	_ = hub1.publish(ctx, &Message{data: []byte("test-04")})
	_ = hub2.publish(ctx, &Message{data: []byte("test-05")})

	time.Sleep(2 * time.Second)
}
