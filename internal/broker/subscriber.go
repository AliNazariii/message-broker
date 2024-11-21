package broker

import (
	"therealbroker/pkg/broker"
)

type Subscriber struct {
	id     int
	stream chan broker.Message

	// isStopped is true on context cancel
	isStopped bool
}

func NewSubscriber(id int) *Subscriber {
	return &Subscriber{
		id,
		make(chan broker.Message, 200),
		false,
	}
}
