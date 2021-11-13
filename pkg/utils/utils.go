package utils

import "os"

// GetEnv takes in the environmental variable key and a fallback
// if the env var is absent then the fallback is returned or else the
// variable will be returned
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
