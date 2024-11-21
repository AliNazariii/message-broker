package broker

import (
	"context"
	"errors"
	"sync"
	"therealbroker/internal/repositories"
	"therealbroker/pkg/broker"
	"time"
)

type Module struct {
	latestSubscriberID int
	subscribers        map[string][]*Subscriber
	subscribersLock    sync.Mutex

	publisherLock sync.Mutex
	latestIDs     map[string]int
	isClosed      bool

	messagesRepo repositories.MessageRepo
}

func NewModule(messagesRepo repositories.MessageRepo) broker.Broker {
	latestIDs, _ := messagesRepo.GetTopicLatestIDs()

	return &Module{
		subscribers:  make(map[string][]*Subscriber),
		isClosed:     false,
		latestIDs:    latestIDs,
		messagesRepo: messagesRepo,
	}
}

func (m *Module) Close() error {
	if m.isClosed == true {
		return nil
	}

	for _, subscribers := range m.subscribers {
		for _, subscriber := range subscribers {
			close(subscriber.stream)
		}
	}

	m.isClosed = true

	return nil
}

func (m *Module) Publish(ctx context.Context, subject string, msg broker.Message) (int, error) {
	if m.isClosed == true {
		return 0, broker.ErrUnavailable
	}

	m.publisherLock.Lock()
	m.latestIDs[subject]++
	id := m.latestIDs[subject]
	m.publisherLock.Unlock()

	var wg sync.WaitGroup

	if msg.Expiration != 0 {
		wg.Add(1)
		go func(id int, subject string, msg broker.Message) {
			defer wg.Done()
			m.messagesRepo.AddMessage(id, subject, msg.Body, int64(msg.Expiration))
		}(id, subject, msg)
	}

	for _, subscriber := range m.subscribers[subject] {
		if !subscriber.isStopped {
			wg.Add(1)
			go func(subscriber *Subscriber) {
				defer wg.Done()
				subscriber.stream <- msg
			}(subscriber)
		}
	}

	wg.Wait()
	return id, nil
}

func (m *Module) Subscribe(ctx context.Context, subject string) (<-chan broker.Message, error) {
	if m.isClosed == true {
		return nil, broker.ErrUnavailable
	}

	m.subscribersLock.Lock()
	m.latestSubscriberID++
	id := m.latestSubscriberID
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
	if m.isClosed == true {
		return broker.Message{}, broker.ErrUnavailable
	}

	message, err := m.messagesRepo.GetByTopicAndID(id, subject)
	if err != nil {
		return broker.Message{}, broker.ErrInvalidID
	}

	if time.Now().After(message.CreatedAt.Add(time.Duration(message.Expiration))) {
		return broker.Message{}, broker.ErrExpiredID
	}

	return broker.Message{
		Body:       message.Body,
		Expiration: time.Duration(message.Expiration),
	}, nil
}
