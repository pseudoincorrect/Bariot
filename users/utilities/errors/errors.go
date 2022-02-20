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

func NewDbError(text string) error {
	return &AppError{"Database error: " + text}
}

func NewPasswordError() error {
	return &AppError{"Incorect password"}
}

func NewAuthError(text string) error {
	return &AppError{"Authentication/Authorization error: " + text}
}
