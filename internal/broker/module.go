package broker

import (
	"context"
	"errors"
	"sync"
	"therealbroker/internal/repositories"
	"therealbroker/pkg/broker"
	"therealbroker/pkg/log"
	"time"
)

type Subscriber struct {
	id     int
	stream chan broker.Message
	// isStopped is true on context cancel
	isStopped bool
}

type Module struct {
	latestSubscriberId int
	subscribers        map[string][]*Subscriber
	subscribersLock    sync.Mutex

	publisherLock sync.Mutex
	latestIds     map[string]int
	isClosed      bool

	log         *log.Logger
	MessageRepo repositories.MessageRepo
}

func NewModule(log *log.Logger, messageRepo repositories.MessageRepo) broker.Broker {
	log.Debugln("NewModule")
	subscribers := make(map[string][]*Subscriber)
	latestIds, _ := messageRepo.GetTopicLatestIds()
	log.Debugln("latestIds", latestIds)
	return &Module{subscribers: subscribers, isClosed: false, latestIds: latestIds, log: log, MessageRepo: messageRepo}
}

func (m *Module) Close() error {
	m.log.Debugln("Close")
	if m.isClosed == true {
		return nil
	}

	for _, subscribers := range m.subscribers {
		for _, subscriber := range subscribers {
			close(subscriber.stream)
		}
	}
	m.log.Debugln("Closed")
	m.isClosed = true
	return nil
}

func (m *Module) Publish(ctx context.Context, subject string, msg broker.Message) (int, error) {
	//m.log.Debugln("Publish")
	if m.isClosed == true {
		return 0, broker.ErrUnavailable
	}

	m.publisherLock.Lock()
	m.latestIds[subject]++
	id := m.latestIds[subject]
	m.publisherLock.Unlock()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(id int, subject string, msg broker.Message) {
		defer wg.Done()
		if msg.Expiration != 0 {
			m.MessageRepo.AddMessage(id, subject, msg.Body, int64(msg.Expiration))
		}
	}(id, subject, msg)

	for _, subscriber := range m.subscribers[subject] {
		wg.Add(1)
		go func(subscriber *Subscriber) {
			if subscriber.isStopped == false {
				subscriber.stream <- msg
			}
			wg.Done()
		}(subscriber)
	}

	wg.Wait()
	return id, nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	m.log.Debugln("Subscribe", subject)
	if m.isClosed == true {
		return nil, broker.ErrUnavailable
	}

	m.subscribersLock.Lock()
	m.latestSubscriberId++
	id := m.latestSubscriberId
	m.subscribersLock.Unlock()

	m.subscribersLock.Lock()
	subscriber := NewSubscriber(id)
	m.subscribers[subject] = append(m.subscribers[subject], subscriber)
	m.subscribersLock.Unlock()

	go func(subscriber *Subscriber) {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.Canceled) {
			subscriber.isStopped = true
		}
	}(subscriber)

	return subscriber.stream, nil
}

func (m *Module) Fetch(ctx context.Context, subject string, id int) (broker.Message, error) {
	m.log.Debugln("Fetch")
	if m.isClosed == true {
		return broker.Message{}, broker.ErrUnavailable
	}

	message, err := m.MessageRepo.GetByTopicAndId(id, subject)
	if err != nil {
		m.log.Debugln("MessageRepo.GetByTopicAndId", err.Error())
		return broker.Message{}, broker.ErrInvalidID
	}

	if time.Now().Before(message.CreatedAt.Add(time.Duration(message.Expiration))) {
		return broker.Message{Body: message.Body, Expiration: time.Duration(message.Expiration)}, nil
	}
	m.log.Debugln("expired", "now:", time.Now(), "createdAt:", message.CreatedAt, "duration:", message.Expiration)
	return broker.Message{}, broker.ErrExpiredID
}
