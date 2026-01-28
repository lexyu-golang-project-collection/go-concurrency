package main

import (
	"context"
	"log"
	"sync"
)

type Hub struct {
	mu          sync.Mutex
	subscribers map[*Subscriber]struct{}
}

func NewHub() *Hub {
	return &Hub{
		subscribers: map[*Subscriber]struct{}{},
	}
}

func (hub *Hub) publish(ctx context.Context, msg *Message) error {
	hub.mu.Lock()
	log.Printf("Hub broadcasting to %d subscribers", len(hub.subscribers))
	for sub := range hub.subscribers {
		sub.publish(ctx, msg)
	}
	hub.mu.Unlock()
	return nil
}

func (hub *Hub) subscribe(ctx context.Context, subscribe *Subscriber) error {
	hub.mu.Lock()
	hub.subscribers[subscribe] = struct{}{}
	hub.mu.Unlock()

	go func() {
		select {
		case <-subscribe.quit:
		case <-ctx.Done():
			hub.mu.Lock()
			delete(hub.subscribers, subscribe)
			hub.mu.Unlock()
		}
	}()

	go subscribe.run(ctx)

	return nil
}

func (hub *Hub) unSubscribe(subscribe *Subscriber) error {
	hub.mu.Lock()
	log.Printf("%+v cancel subscribe", subscribe.name)
	delete(hub.subscribers, subscribe)
	hub.mu.Unlock()
	close(subscribe.quit)
	return nil
}

func (hub *Hub) getUsers() int {
	hub.mu.Lock()
	length := len(hub.subscribers)
	hub.mu.Unlock()
	return length
}
