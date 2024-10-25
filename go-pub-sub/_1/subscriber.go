package main

import (
	"context"
	"log"
	"sync"
)

type Subscriber struct {
	mu sync.Mutex

	name    string
	channel chan *Message
	quit    chan struct{}
}

func (sub *Subscriber) run(ctx context.Context) {
	for {
		select {
		case msg := <-sub.channel:
			log.Println(sub.name, string(msg.data))
		case <-sub.quit:
			return
		case <-ctx.Done():
			return
		}
	}
}

func (sub *Subscriber) publish(ctx context.Context, msg *Message) {
	select {
	case <-ctx.Done():
		return
	case sub.channel <- msg:
	default:
	}
}

func NewSubscriber(name string) *Subscriber {
	return &Subscriber{
		name:    name,
		channel: make(chan *Message, 100),
		quit:    make(chan struct{}),
	}
}
