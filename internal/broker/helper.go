package broker

import (
	"therealbroker/pkg/broker"
	"time"
)

func NewSubscriber(id int) *Subscriber {
	subscriber := make(chan broker.Message, 200)
	return &Subscriber{id, subscriber, false}
}

func CreateBrokerMessage(body []byte, expirationTime int32) broker.Message {
	return broker.Message{Body: string(body), Expiration: time.Duration(expirationTime) * time.Second}
}
