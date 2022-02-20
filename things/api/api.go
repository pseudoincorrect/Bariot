package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pseudoincorrect/bariot/things/models"
	"github.com/pseudoincorrect/bariot/things/service"
	utils "github.com/pseudoincorrect/bariot/things/utilities"
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
	router.Get("/{id}", thingGetEndpoint(s))
	router.Post("/", thingPostEndpoint(s))
	router.Delete("/{id}", thingDeleteEndpoint(s))
	router.Put("/{id}", thingPutEndpoint(s))
}

func startRouter(port string, router *chi.Mux) {
	addr := ":" + port
	http.ListenAndServe(addr, router)
}

func thingGetEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if err := utils.ValidateUuid(id); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		thing, err := s.GetThing(context.Background(), id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if thing == nil {
			http.Error(res, errors.NewThingNotFoundError(id).Error(), http.StatusNotFound)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(thing)
	}
}

type thingPostRequest struct {
	Name   string `json:"Name"`
	Key    string `json:"Key"`
	UserId string `json:"UserId"`
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

func thingPostEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		thingReq := thingPostRequest{}
		err := json.NewDecoder(req.Body).Decode(&thingReq)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = thingReq.validate(); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		thing := models.Thing{
			Key:    thingReq.Key,
			Name:   thingReq.Name,
			UserId: thingReq.UserId,
		}
		savedThing, err := s.SaveThing(context.Background(), &thing)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(savedThing)
	}
}

type thingDeleteResponse struct {
	Id string `json:"Id"`
}

func thingDeleteEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if err := utils.ValidateUuid(id); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		thingId, err := s.DeleteThing(context.Background(), id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(thingDeleteResponse{Id: thingId})
	}
}

type thingPutRequest struct {
	Name   string `json:"Name"`
	Key    string `json:"Key"`
	UserId string `json:"UserId"`
}

func (r *thingPutRequest) validate() error {
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

func thingPutEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if err := utils.ValidateUuid(id); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		thingReq := thingPutRequest{}
		err := json.NewDecoder(req.Body).Decode(&thingReq)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		thing := models.Thing{
			Id:     id,
			Key:    thingReq.Key,
			Name:   thingReq.Name,
			UserId: thingReq.UserId,
		}
		updatedThing, err := s.UpdateThing(context.Background(), &thing)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(updatedThing)
	}
}
