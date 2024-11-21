package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"therealbroker/pkg/config"
	"time"
)

type DB struct {
	Session *gocql.Session
}

func NewDB(conf *config.Cassandra) *DB {
	db := DB{}

	consistency, err := gocql.ParseConsistencyWrapper(conf.Consistency)
	if err != nil {
		consistency = gocql.LocalOne
		logrus.Infof("Error in consistency, so set to %s", consistency)
	}

	cluster := gocql.NewCluster(conf.Hosts...)
	cluster.Keyspace = conf.KeySpace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: conf.Username,
		Password: conf.Password}
	cluster.PageSize = conf.PageSize
	cluster.Port = conf.Port
	cluster.Consistency = consistency
	cluster.Timeout = time.Duration(conf.Timeout) * time.Millisecond
	cluster.WriteCoalesceWaitTime = time.Millisecond
	cluster.NumConns = 10
	if conf.DataCenter != "" {
		cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy(conf.DataCenter))
		cluster.HostFilter = gocql.DataCentreHostFilter(conf.DataCenter)
	}

	db.Session, err = cluster.CreateSession()

	if err != nil {
		logrus.Fatalf("Unable to create cassandra session. %s", err)
	}
	return &db
}
