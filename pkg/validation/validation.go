package validation

import (
	"net/mail"

	"github.com/google/uuid"
	appErr "github.com/pseudoincorrect/bariot/pkg/errors"
)

func ValidateUuid(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return appErr.ErrValidation
	}
	return nil
}

func ValidateEmail(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return appErr.ErrValidation
	}
	return nil
}
