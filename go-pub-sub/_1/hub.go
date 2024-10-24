package main

import (
	"context"
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

func (hub *Hub) subscribe(ctx context.Context, subscribe *Subscriber) error {
	panic("")
}

func (hub *Hub) unSubscribe(ctx context.Context, subscribe *Subscriber) error {
	panic("")
}
