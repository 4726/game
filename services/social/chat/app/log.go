package app

import (
	"os"
	"strings"

	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

var logEntry *NSQLogger

func init() {
	entry := logrus.WithFields(logrus.Fields{
		"app": "chat",
	})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	logEntry = &NSQLogger{entry}
}

// NSQLogger implements a logger compatible with nsq
type NSQLogger struct {
	*logrus.Entry
}

func (n *NSQLogger) Output(calldepth int, s string) error {
	if len(s) < 3 {
		return nil
	}

	if len(s) > 3 {
		level := s[:3]
		msg := strings.TrimSpace(s[3:])
		switch level {
		case nsq.LogLevelInfo.String():
			n.Info(msg)
		case nsq.LogLevelWarning.String():
			n.Warn(msg)
		case nsq.LogLevelError.String():
			n.Error(msg)
		}
	}
	return nil
}
