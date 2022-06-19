package errors

import "errors"

var ErrCreation = errors.New("instance creation error")

var ErrConn = errors.New("connection error")

var ErrWrite = errors.New("write error")

var ErrRead = errors.New("read error")

var ErrGrpc = errors.New("grpc error")

var ErrValidation = errors.New("validation error")

var ErrParsing = errors.New("parsing error")

var ErrNotFound = errors.New("not found error")

var ErrDb = errors.New("database error")

var ErrDbUuid = errors.New("database uuid error")

var ErrDbNotFound = errors.New("database not found error")

var ErrPassword = errors.New("password error")

var ErrAuthn = errors.New("authentication error")

var ErrAuthz = errors.New("authorization error")

var ErrCache = errors.New("cache error")
