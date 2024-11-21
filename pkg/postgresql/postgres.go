package postgresql

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"therealbroker/pkg/config"
)

type DB struct {
	*gorm.DB
}

func NewDB(conf *config.Postgres, models ...interface{}) *DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.DB, conf.Pass,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatalf("Unable to connect to postgres: %v", err)
	}

	if err = setupConnectionPool(db, conf); err != nil {
		logrus.Fatalf("Unable to configure connection pool: %v", err)
	}

	if err = db.AutoMigrate(models...); err != nil {
		logrus.Errorf("Migration problem: %v", err)
	}

	return &DB{DB: db}
}

func setupConnectionPool(db *gorm.DB, conf *config.Postgres) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("unable to access postgres: %w", err)
	}

	sqlDB.SetConnMaxLifetime(conf.ConnMaxLifetime)
	sqlDB.SetMaxIdleConns(conf.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConnections)

	return nil
}
