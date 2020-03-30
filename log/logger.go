package log

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
}

// Info logs info level logs
func Info(args ...interface{}) {
	logger.Info(args...)
}

// Debug logs debug level logs
func Debug(args ...interface{}) {
	logger.Debug(args...)
}

// Error logs error level logs
func Error(args ...interface{}) {
	logger.Error(args...)
}

// SetLogLevel sets the logger level.
func SetLogLevel(level string) {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Panic(err)
	}
	logger.SetLevel(logLevel)
}
