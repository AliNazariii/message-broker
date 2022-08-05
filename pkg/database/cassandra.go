package database

import (
	"github.com/gocql/gocql"
	"therealbroker/pkg/config"
	"therealbroker/pkg/log"
	"time"
)

type CassandraDB struct {
	log     *log.Logger
	Session *gocql.Session
}

func NewCassandraDB(log *log.Logger, conf *config.SectionCassandra) *CassandraDB {
	db := CassandraDB{}

	consistency, err := gocql.ParseConsistencyWrapper(conf.Consistency)
	if err != nil {
		consistency = gocql.LocalOne
		log.Infof("Error in consistency, so set to %s", consistency)
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
		log.Fatalf("Unable to create cassandra session. %s", err)
	}
	return &db
}
