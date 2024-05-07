package mq

import "context"

type MQ interface {
	// send some data to a topic
	Send(topic string, data interface{}) error
	SendWithContext(ctx context.Context, topic string, data interface{}) error

	// or, pass messages back to client
	Receive(topic string) (<-chan []byte, error)
	ReceiveWithContext(ctx context.Context, topic string) (<-chan []byte, error)
}
