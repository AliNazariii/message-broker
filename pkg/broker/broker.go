package broker

import (
	"context"
	"io"
	"time"
)

type Message struct {
	// Identifier for the message, optional.
	// If not provided, the message is not accessible through Fetch().
	// Must be unique for each subject.
	id int
	// Content of the message.
	Body string
	// Duration for which the message remains accessible through Fetch().
	// Set to 0 for messages that do not need to be retained (fire-and-forget).
	Expiration time.Duration
}

// Broker defines a thread-safe interface for message publishing and retrieval.
// Appropriate errors are returned based on errors.go.
type Broker interface {
	io.Closer
	// Publish publishes a message and returns a unique ID.
	// Guarantees message ordering; for instance, if messages A, B, and C are published,
	// subscribers receive them in the same order: A, B, C.
	Publish(ctx context.Context, subject string, msg Message) (int, error)

	// Subscribe returns a channel that receives messages for a given subject.
	// When the context is canceled, message delivery to the subscriber stops.
	// No action occurs on timeout.
	Subscribe(ctx context.Context, subject string) (<-chan Message, error)

	// Fetch returns a previously published message by ID if it has not expired.
	Fetch(ctx context.Context, subject string, id int) (Message, error)
}
