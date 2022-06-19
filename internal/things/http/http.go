package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pseudoincorrect/bariot/internal/things/models"
	"github.com/pseudoincorrect/bariot/internal/things/service"
	e "github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/pkg/validation"
)

type ctxKey int

const userIdKey ctxKey = iota

// InitApi initialize the thing REST api
func InitApi(port string, s service.Things) error {
	router := createRouter()
	createEndpoint(s, router)
	err := startRouter(port, router)
	return err
}

// createRouter create a REST api router with middleware
func createRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}

// createEndpoint create the endpoints for the REST api with authorization
func createEndpoint(s service.Things, r *chi.Mux) {
	// only user can create a thing (associated with user id)
	userOnlyMidGroup := r.Group(nil)
	userOnlyMidGroup.Use(userOnlyMid(s))
	userOnlyMidGroup.Post("/", thingPostEndpoint(s))
	// only user of thing or admin can get delete a thing
	userOfThingOrAdminGroup := r.Group(nil)
	userOfThingOrAdminGroup.Use(userOfThingOrAdminMid(s))
	userOfThingOrAdminGroup.Get("/{id}", thingGetEndpoint(s))
	userOfThingOrAdminGroup.Get("/{id}/token", thingGetTokenEndpoint(s))
	userOfThingOrAdminGroup.Delete("/{id}", thingDeleteEndpoint(s))
	// only a user of a thing can update it
	userOfThingOnlyGroup := r.Group(nil)
	userOfThingOnlyGroup.Use(userOfThingOnlyMid(s))
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
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		thing, err := s.GetThing(context.Background(), id)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if thing == nil {
			e.HandleHttp(res, e.ErrNotFound.Error(), http.StatusNotFound)
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

// validate the thingPostRequest
func (r *thingPostRequest) validate() error {
	if r.Name == "" {
		return e.ErrValidation
	}
	if len(r.Name) > 100 {
		return e.ErrValidation
	}
	if len(r.Name) < 3 {
		return e.ErrValidation
	}
	return nil
}

// thingPostEndpoint create a handler to create a thing
func thingPostEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userId := req.Context().Value(userIdKey).(string)
		thingReq := thingPostRequest{}
		err := json.NewDecoder(req.Body).Decode(&thingReq)
		if err != nil {
			e.HandleHttp(res, e.ErrParsing.Error(), http.StatusBadRequest)
			return
		}
		if err = thingReq.validate(); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		thing := models.Thing{
			Key:    thingReq.Key,
			Name:   thingReq.Name,
			UserId: userId,
		}
		log.Println(thing)
		err = s.SaveThing(context.Background(), &thing)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println(thing)

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(&thing)
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
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		thingId, err := s.DeleteThing(context.Background(), id)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(thingDeleteResponse{Id: thingId})
	}
}

type thingPutRequest struct {
	Name string `json:"Name"`
	Key  string `json:"Key"`
}

// validate the thingPutRequest
func (r *thingPutRequest) validate() error {
	if r.Name == "" {
		return e.ErrValidation
	}
	if len(r.Name) > 100 {
		return e.ErrValidation
	}
	if len(r.Name) < 3 {
		return e.ErrValidation
	}
	return nil
}

// thingPutEndpoint create a handler to update a thing with a thing model
func thingPutEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userId := req.Context().Value(userIdKey).(string)
		id := chi.URLParam(req, "id")
		if err := validation.ValidateUuid(id); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		thingReq := thingPutRequest{}
		err := json.NewDecoder(req.Body).Decode(&thingReq)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = thingReq.validate(); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		thing := models.Thing{
			Id:     id,
			Key:    thingReq.Key,
			Name:   thingReq.Name,
			UserId: userId,
		}
		log.Println("TODO: check thing belong to user")
		err = s.UpdateThing(context.Background(), &thing)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(thing)
	}
}

type thingGetTokenRes struct {
	Token string
}

// thingGetTokenEndpoint create a handler to get a token associated to a thing
func thingGetTokenEndpoint(s service.Things) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userId := req.Context().Value(userIdKey).(string)
		thingId := chi.URLParam(req, "id")
		jwt, err := s.GetThingToken(context.Background(), thingId, userId)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(thingGetTokenRes{Token: jwt})
	}
}

type middlewareFunc func(http.Handler) http.Handler

// userOfThingOrAdmin middleware to check whether the token belong to an admin
// or to the user (ID) of the thing in the request
func userOfThingOrAdminMid(s service.Things) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return userOfThingOrAdmin(s, next)
	}
}

func userOfThingOrAdmin(
	s service.Things, next http.Handler) http.HandlerFunc {
	// middleware logic
	fn := func(res http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		log.Println("token = ", token)
		thingId := chi.URLParam(req, "id")
		if err := validation.ValidateUuid(thingId); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		userId, err := s.UserOfThingOrAdmin(context.Background(), token, thingId)
		if err != nil {
			log.Println(err)
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(res, req.WithContext(context.WithValue(req.Context(), userIdKey, userId)))
	}
	return http.HandlerFunc(fn)
}

// userOfThingOnly middleware to check whether the token in the request belong
// to the user of the thing in the request
func userOfThingOnlyMid(s service.Things) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return userOfThingOnly(s, next)
	}
}

func userOfThingOnly(
	s service.Things, next http.Handler) http.HandlerFunc {
	// middleware logic
	fn := func(res http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		thingId := chi.URLParam(req, "id")
		if err := validation.ValidateUuid(thingId); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		userId, err := s.UserOfThingOnly(context.Background(), token, thingId)
		if err != nil {
			log.Println(err)
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(res, req.WithContext(context.WithValue(req.Context(), userIdKey, userId)))
	}
	return http.HandlerFunc(fn)
}

// userOnlyMid middleware to check whether the token belong to a user, and not an admin
func userOnlyMid(s service.Things) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return userOnly(s, next)
	}
}

func userOnly(
	s service.Things, next http.Handler) http.HandlerFunc {
	// middleware logic
	fn := func(res http.ResponseWriter, req *http.Request) {
		token := req.Header.Get("Authorization")
		userId, err := s.UserOnly(context.Background(), token)
		if err != nil {
			log.Println(err)
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(res,
			req.WithContext(context.WithValue(
				req.Context(),
				userIdKey,
				userId)))
	}

	return http.HandlerFunc(fn)
}
