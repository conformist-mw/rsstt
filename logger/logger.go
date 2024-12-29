package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init(level string) {
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	Log.SetOutput(os.Stdout)

	logLevel := logrus.InfoLevel
	if level, err := logrus.ParseLevel(level); err == nil {
		logLevel = level
	}
	Log.SetLevel(logLevel)
}
