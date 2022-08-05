package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strconv"
	"therealbroker/pkg/config"
	"therealbroker/pkg/log"
	"time"
)

type PostgresDB struct {
	log *log.Logger
	DB  *gorm.DB
}

func NewPostgresDB(log *log.Logger, conf *config.SectionPostgres, models ...interface{}) *PostgresDB {
	var err error
	strConf := "host=" + conf.Host + " port=" + strconv.Itoa(conf.Port) +
		" user=" + conf.User + " dbname=" + conf.DB +
		" password=" + conf.Pass + " sslmode=disable"

	db, err := gorm.Open(postgres.Open(strConf), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to postgres: %v", err)
	}

	pdb, err := db.DB()
	if err != nil {
		log.Fatalf("Unable to access postgres: %v", err)
	}

	pdb.SetConnMaxLifetime(time.Minute * 2)
	pdb.SetMaxIdleConns(20)
	pdb.SetMaxOpenConns(20)
	err = db.AutoMigrate(models...)
	if err != nil {
		log.Errorf("Migration problem is: %v", err)
	}

	return &PostgresDB{log, db}
}
