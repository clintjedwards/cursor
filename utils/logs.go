package utils

import "github.com/sirupsen/logrus"

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
	LogLevelFatal LogLevel = "fatal"
)

// StructuredLog allows the application to log to stdout, json formatted,
//  levels accepted are debug, info, warn, error, and fatal
func StructuredLog(level LogLevel, description string, object interface{}) {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	logger := logrus.WithFields(logrus.Fields{
		"data": object,
	})

	switch level {
	case LogLevelDebug:
		logger.Debugln(description)
	case LogLevelInfo:
		logger.Infoln(description)
	case LogLevelWarn:
		logger.Warnln(description)
	case LogLevelError:
		logger.Errorln(description)
	case LogLevelFatal:
		logger.Fatalln(description)
	default:
		logger.Infoln(description)
	}
}
