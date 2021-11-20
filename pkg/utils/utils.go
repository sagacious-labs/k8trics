package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// GetEnv takes in the environmental variable key and a fallback
// if the env var is absent then the fallback is returned or else the
// variable will be returned
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func SetupLogger() {
	logrus.SetReportCaller(true)

	logLevel := logrus.WarnLevel

	switch GetEnv("K8TRICS_LOG_LEVEL", "warn") {
	case "panic":
		logLevel = logrus.PanicLevel
	case "fatal":
		logLevel = logrus.FatalLevel
	case "error":
		logLevel = logrus.ErrorLevel
	case "warn":
		logLevel = logrus.WarnLevel
	case "info":
		logLevel = logrus.InfoLevel
	case "debug":
		logLevel = logrus.DebugLevel
	case "trace":
		logLevel = logrus.TraceLevel
	}

	logrus.SetLevel(logLevel)
}
