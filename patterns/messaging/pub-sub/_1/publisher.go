package main

import (
	"context"
	"log"
)

type Publisher struct {
	hubs []*Hub
}

func NewPublisher(hubs ...*Hub) *Publisher {
	return &Publisher{
		hubs: hubs,
	}
}

func (p *Publisher) publish(ctx context.Context, msg *Message) error {
	log.Printf("Publisher broadcasting message: %s", string(msg.data))
	for _, hub := range p.hubs {
		if err := hub.publish(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
