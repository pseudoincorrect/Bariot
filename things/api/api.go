package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func StartRouter() {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from Things Service"))
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from Things Service"))
	})

	http.ListenAndServe(":8080", r)
}
