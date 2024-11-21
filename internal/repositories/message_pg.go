package repositories

import (
	"errors"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"therealbroker/pkg/postgresql"
)

// psql -U admin broker_db
// \dt
// docker run --net=host -e DATA_SOURCE_NAME="postgresql://admin:123456@127.0.0.1:5432/broker_db?sslmode=disable" --name=postgres-exporter quay.io/prometheuscommunity/postgres-exporter

/*
CREATE INDEX borker_message_subject_idx
ON pg_messages (subject);
*/

type MessagePostgres struct {
	PID        int    `gorm:"primary_key;auto_increment:true"`
	ID         int    `gorm:"not null"`
	Subject    string `gorm:"not null"`
	Body       string `gorm:"not null"`
	Expiration int64  `gorm:"not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type MessagesPostgres struct {
	db           *postgresql.DB
	messages     []*MessagePostgres
	channels     []chan bool
	messagesLock sync.Mutex
	ticker       *time.Ticker
}

func NewMessagesPostgres(db *postgresql.DB) *MessagesPostgres {
	messages := make([]*MessagePostgres, 0)
	channels := make([]chan bool, 0)

	ticker := time.NewTicker(time.Microsecond * 5)

	mi := &MessagesPostgres{
		db:       db,
		messages: messages,
		ticker:   ticker,
		channels: channels,
	}

	mi.WriteMessages()

	return mi
}

func (mi *MessagesPostgres) GetTopicLatestIDs() (map[string]int, error) {
	latestIds := make(map[string]int)

	rows, err := mi.db.DB.Model(&MessagePostgres{}).Select("subject, max(id) as latestId").Group("subject").Rows()
	if err != nil {
		return latestIds, err
	}

	for rows.Next() {
		var subject string
		var latestId int

		err = rows.Scan(&subject, &latestId)
		if err != nil {
			return latestIds, err
		}

		latestIds[subject] = latestId
	}
	return latestIds, nil
}

func (mi *MessagesPostgres) GetByTopicAndID(id int, subject string) (Message, error) {
	var message MessagePostgres
	err := mi.db.DB.Where("id = ? AND subject = ?", id, subject).First(&message).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return Message{}, err
	}
	return Message{
		ID:         message.ID,
		Subject:    message.Subject,
		Body:       message.Body,
		Expiration: message.Expiration,
		CreatedAt:  message.CreatedAt}, nil
}

func (mi *MessagesPostgres) AddMessage(id int, subject string, body string, expiration int64) {
	message := MessagePostgres{ID: id, Subject: subject, Body: body, Expiration: expiration}

	//mi.db.DB.Create(&message)

	channel := make(chan bool)
	mi.messagesLock.Lock()
	mi.channels = append(mi.channels, channel)
	mi.messages = append(mi.messages, &message)
	mi.messagesLock.Unlock()

	<-channel
	return
}

func (mi *MessagesPostgres) WriteMessages() {
	go func() {
		latencies := make([]time.Duration, 0)
		for {
			select {
			case <-mi.ticker.C:
				if len(mi.messages) == 0 {
					continue
				}

				mi.messagesLock.Lock()

				messages := mi.messages
				mi.messages = make([]*MessagePostgres, 0)

				channels := mi.channels
				mi.channels = make([]chan bool, 0)

				mi.messagesLock.Unlock()
				start := time.Now()
				mi.db.DB.CreateInBatches(&messages, 150)

				latencies = append(latencies, time.Since(start))
				average := time.Duration(0)
				for _, latency := range latencies {
					average += latency
				}
				average = average / time.Duration(len(latencies))
				logrus.Debugln("Wrote", len(messages), "Average latency:", average)

				for _, channel := range channels {
					channel <- true
				}
			}
		}
	}()
}
