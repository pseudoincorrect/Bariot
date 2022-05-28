package errors

import (
	"log"
	"net/http"
)

func Http(res http.ResponseWriter, msg string, code int) {
	log.Println(msg)
	http.Error(res, msg, code)
}
