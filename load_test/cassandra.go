package main

import (
	"go.uber.org/atomic"
	"math/rand"
	"sync"
	"therealbroker/internal/repositories"
	"therealbroker/pkg/config"
	"therealbroker/pkg/database"
	"therealbroker/pkg/log"
	"time"
)

const CassVUs = 10000
const CassREQUESTS = 10000

var subjects = []rune("abcdefghijklmnopqrstuvwxyz")

func main() {
	logger := log.NewLog("debug")
	conf := config.New("config.yaml", "broker")
	brokerCassandra := database.NewCassandraDB(logger, &conf.Cassandra)
	messageRepo := repositories.NewCasMessageRepo(brokerCassandra, logger, &conf.Cassandra)
	latestIds, _ := messageRepo.GetTopicLatestIds()
	var l sync.Mutex

	overallLatency := 0
	var count atomic.Int64
	overallStart := time.Now()
	for i := 0; i < CassVUs; i++ {
		go func(i int) {
			for j := 0; j < CassREQUESTS; j++ {
				l.Lock()
				subject := string(subjects[rand.Intn(len(subjects))])
				latestIds[subject]++
				id := latestIds[subject]
				l.Unlock()

				start := time.Now()
				messageRepo.AddMessage(id, subject, "asasas", 200000)
				overallLatency += int(time.Since(start).Milliseconds())
				count.Inc()
				logger.Debugln("Throughput", float64(count.Load())/time.Since(overallStart).Seconds(), count.Load())
			}
		}(i)
	}
	<-time.After(time.Minute * 10)
}
