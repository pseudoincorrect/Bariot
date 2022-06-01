package errors

import (
	"log"
	"net/http"
)

func Http(res http.ResponseWriter, msg string, code int) {
	log.Println("HTTP ERROR", msg)
	http.Error(res, msg, code)
}
