package config

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"therealbroker/pkg/config"
)

type Config struct {
	Core       Core
	Postgres   config.Postgres
	Prometheus config.Prometheus
	Log        config.Log
	Grpc       config.Grpc
	Jaeger     config.Jaeger
	Cassandra  config.Cassandra
}

type Core struct {
	ServiceName string
}

func New(serviceName string) *Config {
	var conf Config

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	setDefaults(serviceName)

	if err := viper.Unmarshal(&conf); err != nil {
		panic(errors.WithStack(err))
	}

	return &conf
}

func setDefaults(serviceName string) {
	if serviceName == "" {
		serviceName = "Unknown"
	}
	viper.SetDefault("Core.ServiceName", serviceName)

	config.SetDefaults()
}
