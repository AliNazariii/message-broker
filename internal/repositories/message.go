package repositories

import "time"

type Message struct {
	ID         int
	Subject    string
	Body       string
	Expiration int64
	CreatedAt  time.Time
}

type MessageRepo interface {
	GetTopicLatestIDs() (map[string]int, error)
	GetByTopicAndID(id int, subject string) (Message, error)
	AddMessage(id int, subject string, body string, expiration int64)
}
