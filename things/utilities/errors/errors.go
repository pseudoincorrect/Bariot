package errors

type error interface {
	Error() string
}

type validationError struct {
	msg string
}

func (e *validationError) Error() string {
	return e.msg
}

func NewValidationError(text string) error {
	return &validationError{text}
}
