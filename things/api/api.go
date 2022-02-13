package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pseudoincorrect/bariot/things/service"
	"github.com/pseudoincorrect/bariot/things/utilities/errors"
)

func InitApi(port string, s service.Things) {
	router := createRouter()
	createEndpoint(s, router)
	startRouter(port, router)
}

func createRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}

func createEndpoint(s service.Things, router *chi.Mux) {

	router.Get("/", thingGetEndpoint())
	router.Post("/", thingPostEndpoint)
	router.Delete("/", thingDeleteEndpoint())
	router.Put("/", thingPutEndpoint())
}

func startRouter(port string, router *chi.Mux) {
	addr := ":" + port
	http.ListenAndServe(addr, router)
}

func thingGetEndpoint() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Get things endpoint"))
	}
}

type thingPostRequest struct {
	Name string `json:"name"`
}

func (r *thingPostRequest) validate() error {
	if r.Name == "" {
		return errors.NewValidationError("Thing name is required")
	}
	if len(r.Name) > 100 {
		return errors.NewValidationError("Thing name is too long")
	}
	if len(r.Name) < 3 {
		return errors.NewValidationError("Thing name is too short")
	}
	return nil
}

type thingPostResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func thingPostEndpoint(res http.ResponseWriter, req *http.Request) {
	thing := thingPostRequest{}

	err := json.NewDecoder(req.Body).Decode(&thing)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(thing)

	if err = thing.validate(); err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(thingPostResponse{Name: thing.Name, Id: "123456789"})
}

func thingDeleteEndpoint() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Delete things endpoint"))
	}
}

func thingPutEndpoint() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Update things endpoint"))
	}
}
