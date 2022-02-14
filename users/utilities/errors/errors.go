package errors

type error interface {
	Error() string
}

type AppError struct {
	msg string
}

func (e *AppError) Error() string {
	return e.msg
}

func NewValidationError(text string) error {
	return &AppError{text}
}

func NewUserNotFoundError(id string) error {
	return &AppError{"User " + id + " not found"}
}
