package errors

import "errors"

var ErrCreation = errors.New("instance creation error")

var ErrConnection = errors.New("connection error")

var ErrValidation = errors.New("validation error")

var ErrThingNotFound = errors.New("thing not found error")

var ErrUserNotFound = errors.New("user not found error")

var ErrDb = errors.New("database error")

var ErrPassword = errors.New("password error")

var ErrAuthentication = errors.New("authentication error")

var ErrAuthorization = errors.New("authorization error")