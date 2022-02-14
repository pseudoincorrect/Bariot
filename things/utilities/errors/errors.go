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

func NewThingNotFoundError(id string) error {
	return &AppError{"Thing " + id + " not found"}
}
