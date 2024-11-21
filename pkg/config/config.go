package config

import (
	"github.com/spf13/viper"
	"time"
)

type Core struct {
	ServiceName string
}

type Prometheus struct {
	Port string
}

type Log struct {
	Level string
}

type Grpc struct {
	Address string
}

type Graylog struct {
	Level       string
	Host        string
	Port        int
	Facility    string
	Compression bool
}

type Jaeger struct {
	Address string
}

type Cassandra struct {
	Hosts         []string
	Port          int
	Username      string
	Password      string
	KeySpace      string
	Consistency   string
	PageSize      int
	Timeout       int64
	DataCenter    string
	PartitionSize int32
}

type Postgres struct {
	Host               string
	Port               int
	DB                 string
	User               string
	Pass               string
	BatchCount         int
	MaxIdleConnections int
	MaxOpenConnections int
	ConnMaxLifetime    time.Duration
}

func SetDefaults() {
	viper.SetDefault("Postgres.Host", "localhost")
	viper.SetDefault("Postgres.Port", 5432)
	viper.SetDefault("Postgres.DB", "defaultdb")
	viper.SetDefault("Postgres.User", "user")
	viper.SetDefault("Postgres.Pass", "password")
	viper.SetDefault("Postgres.BatchCount", 5)
	viper.SetDefault("Postgres.MaxIdleConnections", 10)
	viper.SetDefault("Postgres.MaxOpenConnections", 50)
	viper.SetDefault("Postgres.ConnMaxLifetime", time.Minute*5)

	viper.SetDefault("Prometheus.Port", 9000)

	viper.SetDefault("Log.Level", "debug")

	viper.SetDefault("Grpc.Address", "127.0.0.1:8888")

	viper.SetDefault("Jaeger.Address", "localhost:6831")

	viper.SetDefault("Cassandra.Hosts", []string{"127.0.0.1"})
	viper.SetDefault("Cassandra.Port", 9042)
	viper.SetDefault("Cassandra.Username", "cassandra")
	viper.SetDefault("Cassandra.Password", "cassandra")
	viper.SetDefault("Cassandra.KeySpace", "nazari")
	viper.SetDefault("Cassandra.Consistency", "LOCAL_ONE")
	viper.SetDefault("Cassandra.PageSize", 5000)
	viper.SetDefault("Cassandra.Timeout", 16000)
	viper.SetDefault("Cassandra.DataCenter", "dc1")
	viper.SetDefault("Cassandra.PartitionSize", 10)
}
