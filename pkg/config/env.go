package config

import (
	"errors"
	"os"
)

var EnvtVarNotFoundError = errors.New("Missing environment variable")

// GetEnv return the current `key` value or `fallback`.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// RequireEnv returns the current `key` value or an error
func RequireEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", EnvtVarNotFoundError
}
