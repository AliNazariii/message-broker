package logger

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
)

func Configure(level string) {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	parseLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Error(errors.New("failed to parse level"))
	}

	logrus.SetLevel(parseLevel)
	logrus.SetOutput(os.Stdout)
}
