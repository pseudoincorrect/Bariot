package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	appErr "github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/pkg/validation"
	"github.com/pseudoincorrect/bariot/things/models"
	"github.com/pseudoincorrect/bariot/things/service"
)

type ctxKey int

const userIdKey ctxKey = iota

func InitApi(port string, s service.Things) error {
	router := createRouter()
	createEndpoint(s, router)
	err := startRouter(port, router)
	return err
}

func createRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}

// createEndpoint Create all things related endpoints
func createEndpoint(s service.Things, r *chi.Mux) {
	// only user can create a thing (associated with user id)
	userOnlyGroup := r.Group(nil)
	userOnlyGroup.Use(userOnly(s))
	userOnlyGroup.Post("/", thingPostEndpoint(s))
	// only user of thing or admin can get delete a thing
	userOfThingOrAdminGroup := r.Group(nil)
	userOfThingOrAdminGroup.Use(userOfThingOrAdmin(s))
	userOfThingOrAdminGroup.Get("/{id}", thingGetEndpoint(s))
	userOfThingOrAdminGroup.Get("/{id}/token", thingGetTokenEndpoint(s))
	userOfThingOrAdminGroup.Delete("/{id}", thingDeleteEndpoint(s))
	// only a user of a thing can update it
	userOfThingOnlyGroup := r.Group(nil)
	userOfThingOnlyGroup.Use(userOfThingOnly(s))
	userOfThingOnlyGroup.Put("/{id}", thingPutEndpoint(s))
}

// startRouter start the chi http router
func startRouter(port string, r *chi.Mux) error {
	addr := ":" + port
	err := http.ListenAndServe(addr, r)
	return err
}

// thingGetEndpoint create a handler to get a Thing by ID
func thingGetEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if err := validation.ValidateUuid(id); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		thing, err := s.GetThing(context.Background(), id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if thing == nil {
			http.Error(res, "thing not found", http.StatusNotFound)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(thing)
	}
}

type thingPostRequest struct {
	Name string `json:"Name"`
	Key  string `json:"Key"`
}

func (r *thingPostRequest) validate() error {
	if r.Name == "" {
		return appErr.ErrValidation
	}
	if len(r.Name) > 100 {
		return appErr.ErrValidation
	}
	if len(r.Name) < 3 {
		return appErr.ErrValidation
	}
	return nil
}

// thingPostEndpoint create a handler to create a thing
func thingPostEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userId := req.Context().Value(userIdKey).(string)
		log.Println("user id", userId)
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
			UserId: userId,
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

// thingDeleteEndpoint create a handler to delete a thing by ID
func thingDeleteEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if err := validation.ValidateUuid(id); err != nil {
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
	Name string `json:"Name"`
	Key  string `json:"Key"`
}

func (r *thingPutRequest) validate() error {
	if r.Name == "" {
		return appErr.ErrValidation
	}
	if len(r.Name) > 100 {
		return appErr.ErrValidation
	}
	if len(r.Name) < 3 {
		return appErr.ErrValidation
	}
	return nil
}

// thingPutEndpoint create a handler to update a thing with a thing model
func thingPutEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userId := req.Context().Value(userIdKey).(string)
		id := chi.URLParam(req, "id")
		if err := validation.ValidateUuid(id); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		thingReq := thingPutRequest{}
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
			Id:     id,
			Key:    thingReq.Key,
			Name:   thingReq.Name,
			UserId: userId,
		}
		log.Println("TODO: check thing belong to user")
		updatedThing, err := s.UpdateThing(context.Background(), &thing)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(updatedThing)
	}
}

type thingGetTokenRes struct {
	Jwt string
}

// thingGetTokenEndpoint create a handler to get a token associated to a thing
func thingGetTokenEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userId := req.Context().Value(userIdKey).(string)
		thingId := chi.URLParam(req, "id")
		jwt, err := s.GetThingToken(context.Background(), thingId, userId)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(thingGetTokenRes{Jwt: jwt})
	}
}

type middlewareFunc func(http.Handler) http.Handler

// userOfThingOrAdmin middleware to check whether the token belong to an admin
// or to the user (ID) of the thing in the request
func userOfThingOrAdmin(s service.Things) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("Authorization")
			thingId := chi.URLParam(req, "id")
			if err := validation.ValidateUuid(thingId); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			userId, err := s.UserOfThingOrAdmin(context.Background(), token, thingId)
			if err != nil {
				log.Println(err)
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(res, req.WithContext(context.WithValue(req.Context(), userIdKey, userId)))
		})
	}
}

// userOfThingOnly middleware to check whether the token in the request belong
// to the user of the thing in the request
func userOfThingOnly(s service.Things) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("Authorization")
			thingId := chi.URLParam(req, "id")
			if err := validation.ValidateUuid(thingId); err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}
			userId, err := s.UserOfThingOnly(context.Background(), token, thingId)
			if err != nil {
				log.Println(err)
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(res, req.WithContext(context.WithValue(req.Context(), userIdKey, userId)))
		})
	}
}

// userOnly middleware to check whether the token belong to a user, and not an admin
func userOnly(s service.Things) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("Authorization")
			userId, err := s.UserOnly(context.Background(), token)
			if err != nil {
				log.Println(err)
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(res, req.WithContext(context.WithValue(req.Context(), userIdKey, userId)))
		})
	}
}
