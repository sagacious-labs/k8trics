package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
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

// SetupLogger reads environmental variable and sets up logrus
// logger log level
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

// TrimPodTemplateHash takes in a pod and tries to remove the pod template
// hash from the pod name. If no pod template hash is found then pod name
// is returned without any trimming
func TrimPodTemplateHash(pod *v1.Pod) string {
	labels := pod.GetLabels()
	podTemplateHash, ok := labels["pod-template-hash"]
	if !ok {
		return pod.GetName()
	}

	return strings.TrimSuffix(pod.GetName(), fmt.Sprintf("-%s", podTemplateHash))
}
