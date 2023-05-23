package logs

import (
	"github.com/sirupsen/logrus"
)

func Init(debugMode bool) {
	SetFormatter(logrus.StandardLogger())

	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func SetFormatter(logger *logrus.Logger) {
	logger.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
	})
}
