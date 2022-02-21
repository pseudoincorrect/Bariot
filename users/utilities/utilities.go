package utilities

import (
	"log"
	"os"

	"github.com/google/uuid"
	appErr "github.com/pseudoincorrect/bariot/users/utilities/errors"
	"golang.org/x/crypto/bcrypt"
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
		log.Println("Environment variable", key, "is not set")
		log.Println("Please set it and try again")
		panic("Environment variable " + key + " is not set and")
	}
	return value
}

func ValidateUuid(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return appErr.ErrValidation
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
