package main

import (
	"context"
	"log"
)

type Subscriber struct {
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
		log.Printf("Subscriber %s received message", sub.name)
	default:
		log.Printf("Subscriber %s channel full, message dropped", sub.name)
	}
}

func NewSubscriber(name string) *Subscriber {
	return &Subscriber{
		name:    name,
		channel: make(chan *Message, 100),
		quit:    make(chan struct{}),
	}
}
