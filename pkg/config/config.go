package config

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
)

const defaultConf = `
core:
  service_name: "Unknown"
postgres:
  host: ""
  port: 5432
  db: ""
  user: ""
  pass: ""
  batch_count: 5
prometheus:
  port: 8080
log:
  level: debug
grpc:
  address: "127.0.0.1:5050"
jaeger:
  address: "localhost:6831"
cassandra:
  hosts: [ "192.168.70.2" ]
  datacenter: "dc1"
  port: 9042
  user: ""
  password: ""
  keyspace: "nazari"
  consistency: "LOCAL_ONE"
  pagesize: 5000
  timeout: 16000
  partition_size: 10
`

type Config struct {
	defaultConf []byte
	configPath  string
	serviceName string

	Core       SectionCore       `yaml:"core"`
	Postgres   SectionPostgres   `yaml:"postgres"`
	Prometheus SectionPrometheus `yaml:"prometheus"`
	Log        SectionLog        `yaml:"log"`
	Grpc       SectionGrpc       `yaml:"grpc"`
	Jaeger     SectionJaeger     `yaml:"jaeger"`
	Cassandra  SectionCassandra  `yaml:"cassandra"`
}

type SectionCore struct {
	ServiceName string `mapstructure:"service_name"`
}

type SectionPostgres struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	DB         string `yaml:"db"`
	User       string `yaml:"user"`
	Pass       string `yaml:"pass"`
	BatchCount int    `mapstructure:"batch_count"`
}

type SectionPrometheus struct {
	Port string `yaml:"port"`
}

type SectionLog struct {
	Level string `yaml:"level"`
}

type SectionGrpc struct {
	Address string `yaml:"address"`
}

type SectionGraylog struct {
	Level       string `yaml:"level"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Facility    string `yaml:"facility"`
	Compression bool   `yaml:"compression"`
}

type SectionJaeger struct {
	Address string `yaml:"address"`
}

type SectionCassandra struct {
	Hosts         []string `yaml:"hosts"`
	Port          int      `yaml:"port"`
	Username      string   `yaml:"user"`
	Password      string   `yaml:"password"`
	KeySpace      string   `yaml:"keyspace"`
	Consistency   string   `yaml:"consistency"`
	PageSize      int      `yaml:"pagesize"`
	Timeout       int64    `yaml:"timeout"`
	DataCenter    string   `yaml:"datacenter"`
	PartitionSize int32    `yaml:"partition_size"`
}

func (c *Config) configureViper() {
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()                                            // read in environment variables that match
	viper.SetEnvPrefix(strings.ReplaceAll(c.serviceName, "-", "_")) // will be uppercase automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.AddConfigPath("/etc/" + c.serviceName + "/")
	viper.AddConfigPath("$HOME/." + c.serviceName)
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
}

// LoadConf load config from file and read in environment variables that match
func (c *Config) loadConf() error {
	c.configureViper()

	if err := c.readConf(); err != nil {
		return err
	}

	err := viper.Unmarshal(&c)
	if err != nil {
		fmt.Println("unable to decode into config struct, ", err)
		return err
	}

	return nil
}

func readConfFromFile(confPath string) error {
	if confPath != "" {
		content, err := ioutil.ReadFile(confPath)

		if err != nil {
			log.Errorf("File does not exist : %s", confPath)
			return err
		}

		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return err
		}
	} else {
		// If a config file is found, read it in.
		if err := viper.MergeInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			fmt.Println("Config file not found.")
		}
	}
	return nil
}

func (c *Config) readConf() error {
	// load default config
	if err := viper.ReadConfig(bytes.NewBuffer(c.defaultConf)); err != nil {
		return err
	}
	if err := readConfFromFile(c.configPath); err != nil {
		return err
	}
	return nil
}

func New(path string, serviceName string) *Config {
	conf := Config{
		defaultConf: []byte(defaultConf),
		configPath:  path,
		serviceName: serviceName,
	}
	err := conf.loadConf()
	if err != nil {
		log.Fatalf("Load yaml config file error: '%v'", err)
		return nil
	}
	return &conf
}
