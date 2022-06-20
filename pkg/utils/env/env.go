package env

import (
	"os"

	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
)

/// GetEnv returns the value of the environment variable named by the key.
/// If the environment variable is present but empty, GetEnv returns the
/// specified default (Falback Fb) value.
func GetEnvFb(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

/// GetEnv returns the value of the environment variable named by the key.
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Error("Environment variable", key, "is not set")
		logger.Error("Please set it and try again")
		panic("Environment variable " + key + " is not set and")
	}
	return value
}
