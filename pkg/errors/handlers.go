package errors

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pseudoincorrect/bariot/pkg/utils/debug"
)

type AppError interface {
	Error() string
}

type appError struct {
	err    error
	origin error
	msg    string
}

var _ AppError = (*appError)(nil)

func New(err error, origin error, msg string) appError {
	return appError{
		err:    err,
		origin: origin,
		msg:    msg,
	}
}

func (err *appError) Error() string {
	var s string
	if err.msg == "" {
		s = fmt.Sprintf("%s | %s", err.err, err.origin)
	} else {
		s = fmt.Sprintf("%s | %s | %s", err.err, err.msg, err.origin)
	}
	return s
}

func Handle(err error, origin error, msg string) error {
	er := New(err, origin, msg)
	debug.LogWithDepth(3, "ERROR", er.Error())
	return &er
}

func HandleFatal(err error, origin error, msg string) {
	er := New(err, origin, msg)
	debug.LogWithDepth(3, "FATAL ERROR", er.Error())
	log.Fatal("[FATAL ERROR]", er.Error())
}

func HandleHttp(res http.ResponseWriter, msg string, code int) {
	debug.LogWithDepth(3, "HTTP ERROR", msg)
	http.Error(res, msg, code)
}
