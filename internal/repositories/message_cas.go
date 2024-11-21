package repositories

import (
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"

	"therealbroker/pkg/cassandra"
	"therealbroker/pkg/config"
)

// https://www.guru99.com/cassandra-tutorial.html
// brew install java
// brew install cassandra
// https://stackoverflow.com/a/69486477
// cassandra -f
// cqlsh (-u <username> -p <password>) ip
// https://www.cloudwalker.io/2020/05/17/monitoring-cassandra-with-prometheus/

/*
CREATE KEYSPACE nazari_broker
  WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
*/

/*
CREATE TABLE nazari_broker.message (
	id int,
	subject text,
	body text,
	expiration bigint,
	created_at timestamp,
	PRIMARY KEY ((subject), id)
);
*/

type CasMessageImpl struct {
	db           *cassandra.DB
	ticker       *time.Ticker
	messages     []*Message
	channels     []chan bool
	messagesLock sync.Mutex
}

func NewCasMessageRepo(db *cassandra.DB, conf *config.Cassandra) *CasMessageImpl {
	keySpace, _ := db.Session.KeyspaceMetadata(conf.KeySpace)
	logrus.Debugln("Tables:", keySpace.Tables)
	if _, exists := keySpace.Tables["message"]; exists != true {
		err := db.Session.Query("CREATE TABLE message (" +
			"id int, subject text, body text, expiration bigint, created_at timestamp, " +
			"PRIMARY KEY ((subject), id))").Exec()
		if err != nil {
			logrus.Fatal(err)
		}
	}

	messages := make([]*Message, 0)
	channels := make([]chan bool, 0)
	ticker := time.NewTicker(time.Microsecond * 5)
	mi := &CasMessageImpl{
		db:       db,
		ticker:   ticker,
		messages: messages,
		channels: channels,
	}
	mi.WriteMessages()
	return mi
}

func (mi *CasMessageImpl) GetTopicLatestIDs() (map[string]int, error) {
	subjects := make([]string, 0)
	latestIds := make(map[string]int)

	var subject string
	iter := mi.db.Session.Query("SELECT distinct subject FROM message").Iter()
	for iter.Scan(&subject) {
		subjects = append(subjects, subject)
	}

	for _, subject := range subjects {
		var id int
		err := mi.db.Session.Query("SELECT MAX(id) FROM message WHERE subject = ?", subject).Scan(&id)
		if err != nil {
			logrus.Fatal(err)
		}
		latestIds[subject] = id
	}
	return latestIds, nil
}

func (mi *CasMessageImpl) GetByTopicAndID(id int, subject string) (Message, error) {
	var body string
	var expiration int64
	var createdAt time.Time
	err := mi.db.Session.Query(
		"SELECT body, expiration, created_at FROM message "+
			"WHERE "+
			"id="+strconv.Itoa(id)+" AND "+
			"subject='"+subject+"'").Scan(&body, &expiration, &createdAt)
	message := Message{
		ID:         id,
		Subject:    subject,
		Body:       body,
		Expiration: expiration,
		CreatedAt:  createdAt,
	}
	logrus.Debugln(message)
	if errors.Is(err, gocql.ErrNotFound) {
		return Message{}, err
	}
	return message, nil
}

func (mi *CasMessageImpl) AddMessage(id int, subject string, body string, expiration int64) {
	//err := mi.db.Session.Query("INSERT INTO message (id, subject, body, expiration, created_at) " +
	//	"VALUES (" +
	//	"" + strconv.Itoa(id) + ", " +
	//	"'" + subject + "', " +
	//	"'" + body + "', " +
	//	strconv.FormatInt(expiration, 10) + ", " +
	//	"toTimestamp(now()))").Exec()
	//if err != nil {
	//	logrus.Debugln(err)
	//}

	message := Message{ID: id, Subject: subject, Body: body, Expiration: expiration}
	channel := make(chan bool)
	mi.messagesLock.Lock()
	mi.channels = append(mi.channels, channel)
	mi.messages = append(mi.messages, &message)
	mi.messagesLock.Unlock()

	<-channel
	return
}

func (mi *CasMessageImpl) WriteMessages() {
	stmt := "INSERT INTO message (id, subject, body, expiration, created_at) VALUES (?, ?, ?, ?, toTimestamp(now()))"
	var requestCount atomic.Int64
	go func() {
		latencies := make([]time.Duration, 0)
		overallStart := time.Now()
		for {
			select {
			case <-mi.ticker.C:
				if len(mi.messages) == 0 {
					continue
				}

				mi.messagesLock.Lock()

				messages := mi.messages
				mi.messages = make([]*Message, 0)

				channels := mi.channels
				mi.channels = make([]chan bool, 0)

				mi.messagesLock.Unlock()

				start := time.Now()

				batch := mi.db.Session.NewBatch(gocql.UnloggedBatch)
				for i, message := range messages {
					batch.Query(stmt, message.ID, message.Subject, message.Body, message.Expiration)
					if i != 0 && i%500 == 0 {
						requestCount.Inc()
						err := mi.db.Session.ExecuteBatch(batch)
						if err != nil {
							logrus.Errorln("ExecuteBatch", err)
						}
						batch = mi.db.Session.NewBatch(gocql.UnloggedBatch)
					}
				}
				if batch.Size() > 0 {
					requestCount.Inc()
					err := mi.db.Session.ExecuteBatch(batch)
					if err != nil {
						logrus.Errorln("ExecuteBatch", err)
					}
				}

				latencies = append(latencies, time.Since(start))
				average := time.Duration(0)
				for _, latency := range latencies {
					average += latency
				}
				average = average / time.Duration(len(latencies))
				logrus.Debugln("Wrote", len(messages), "Average latency:", average, "Request count:", requestCount.Load(), "Total time:", time.Since(overallStart), "DB Request/sec:", float64(requestCount.Load())/time.Since(overallStart).Seconds())

				for _, channel := range channels {
					channel <- true
				}
			}
		}
	}()
}
