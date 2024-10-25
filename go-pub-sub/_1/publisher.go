package main

import "context"

type Publisher struct {
}

func (publisher *Publisher) publish(ctx context.Context, msg *Message) error {
	panic("")
}
