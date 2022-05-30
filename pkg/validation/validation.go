package validation

import (
	"net/mail"

	"github.com/google/uuid"
	appErr "github.com/pseudoincorrect/bariot/pkg/errors"
)

// ValidateUUID validates a UUID string
func ValidateUuid(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return appErr.ErrValidation
	}
	return nil
}

// ValidateEmail validates an email string
func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return appErr.ErrValidation
	}
	return nil
}
