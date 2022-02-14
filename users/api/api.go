package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pseudoincorrect/bariot/users/models"
	"github.com/pseudoincorrect/bariot/users/service"
	utils "github.com/pseudoincorrect/bariot/users/utilities"
	"github.com/pseudoincorrect/bariot/users/utilities/errors"
)

func InitApi(port string, s service.Users) {
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

func createEndpoint(s service.Users, router *chi.Mux) {
	router.Get("/{id}", userGetEndpoint(s))
	router.Post("/", userPostEndpoint(s))
	router.Delete("/{id}", userDeleteEndpoint(s))
	router.Put("/{id}", userPutEndpoint(s))
}

func startRouter(port string, router *chi.Mux) {
	addr := ":" + port
	http.ListenAndServe(addr, router)
}

func userGetEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		if err := utils.ValidateUuid(id); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := s.GetUser(context.Background(), id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(res, errors.NewUserNotFoundError(id).Error(), http.StatusNotFound)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(user)
	}
}

type userPostRequest struct {
	Email    string `json:"Email"`
	FullName string `json:"FullName"`
}

func (r *userPostRequest) validate() error {
	if r.FullName == "" || len(r.FullName) > 100 || len(r.FullName) < 3 {
		return errors.NewValidationError("Invalid user name")
	}
	if r.Email == "" || len(r.Email) > 100 || len(r.Email) < 3 {
		return errors.NewValidationError("Invalid user name")
	}

	return nil
}

func userPostEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userReq := userPostRequest{}

		err := json.NewDecoder(req.Body).Decode(&userReq)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err = userReq.validate(); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		user := models.User{
			Email:    userReq.Email,
			FullName: userReq.FullName,
		}

		savedUser, err := s.SaveUser(context.Background(), &user)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(savedUser)
	}
}

type userDeleteResponse struct {
	Id string `json:"Id"`
}

func userDeleteEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		if err := utils.ValidateUuid(id); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := s.DeleteUser(context.Background(), id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(res).Encode(userDeleteResponse{Id: userId})
	}
}

type userPutRequest struct {
	Email    string `json:"Email"`
	FullName string `json:"FullName"`
}

func (r *userPutRequest) validate() error {
	if r.FullName == "" || len(r.FullName) > 100 || len(r.FullName) < 3 {
		return errors.NewValidationError("Invalid user name")
	}
	if r.Email == "" || len(r.Email) > 100 || len(r.Email) < 3 {
		return errors.NewValidationError("Invalid user email")
	}

	return nil
}

func userPutEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")

		if err := utils.ValidateUuid(id); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		userReq := userPutRequest{}

		err := json.NewDecoder(req.Body).Decode(&userReq)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		if err = userReq.validate(); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		user := models.User{
			Id:       id,
			Email:    userReq.Email,
			FullName: userReq.FullName,
		}

		updatedUser, err := s.UpdateUser(context.Background(), &user)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(updatedUser)
	}
}
