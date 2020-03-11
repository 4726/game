package app

import (
	"os"

	"github.com/sirupsen/logrus"
)

var logEntry *logrus.Entry

func init() {
	entry := logrus.WithFields(logrus.Fields{
		"app": "5v5fps",
	})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	logEntry = entry
}
