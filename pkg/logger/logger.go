package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func Init(level, format string) {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)

	// Set log level
	switch level {
	case "debug":
		Log.SetLevel(logrus.DebugLevel)
	case "info":
		Log.SetLevel(logrus.InfoLevel)
	case "warn":
		Log.SetLevel(logrus.WarnLevel)
	case "error":
		Log.SetLevel(logrus.ErrorLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	if format == "json" {
		Log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}

func Debug(args ...interface{}) {
	Log.Debug(args...)
}

func Info(args ...interface{}) {
	Log.Info(args...)
}

func Warn(args ...interface{}) {
	Log.Warn(args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}
