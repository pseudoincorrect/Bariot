package errors

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

const debug = true

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
	printCallerInfo()
	er := New(err, origin, msg)
	log.Println("[ERROR]", er.Error())
	return &er
}

func HandleFatal(err error, origin error, msg string) {
	printCallerInfo()
	er := New(err, origin, msg)
	log.Fatal("[FATAL ERROR]", er.Error())
}

func HandleHttp(res http.ResponseWriter, msg string, code int) {
	log.Println("HTTP ERROR", msg)
	http.Error(res, msg, code)
}

func printCallerInfo() {
	if debug {
		_, file, no, ok := runtime.Caller(2)
		if ok {
			splits := strings.Split(file, "/")
			fileName := splits[len(splits)-1]
			log.Println("[ERROR FROM]", fileName, "[LINE]", no)
		}
	}
}
