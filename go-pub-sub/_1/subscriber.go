package main

import (
	"context"
	"sync"
)

type Subscriber struct {
	mu sync.Mutex

	name    string
	channel chan *Message
	quit    chan struct{}
}

func (sub *Subscriber) run(ctx context.Context) {

}

func (s *Subscriber) publish(ctx context.Context, msg *Message) {

}

func NewSubscriber(name string) *Subscriber {
	return &Subscriber{
		name:    name,
		channel: make(chan *Message, 100),
		quit:    make(chan struct{}),
	}
}
