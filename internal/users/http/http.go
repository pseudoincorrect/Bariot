package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pseudoincorrect/bariot/internal/users/hash"
	"github.com/pseudoincorrect/bariot/internal/users/models"
	"github.com/pseudoincorrect/bariot/internal/users/service"
	e "github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/pkg/validation"
)

// InitApi initializes the REST api
func InitApi(port string, s service.Users) error {
	router := createRouter()
	createEndpoint(s, router)
	err := startRouter(port, router)
	return err
}

// createRouter creates the router for the users api with middleware
func createRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}

// createEndpoint creates the endpoints for the users api with authorization
func createEndpoint(s service.Users, r *chi.Mux) {
	// only admins can manage users
	adminGroup := r.Group(nil)
	adminGroup.Use(AdminOnly(s))
	adminGroup.Get("/{id}", userGetEndpoint(s))
	adminGroup.Get("/email/{email}", userGetEmailEndpoint(s))
	adminGroup.Post("/", userPostEndpoint(s))
	adminGroup.Delete("/{id}", userDeleteEndpoint(s))
	adminGroup.Put("/{id}", userPutEndpoint(s))
	r.Post("/login", loginUserEndpoint(s))
	r.Post("/login/admin", loginAdminEndpoint(s))
}

// startRouter starts the HTTP server
func startRouter(port string, r *chi.Mux) error {
	addr := ":" + port
	err := http.ListenAndServe(addr, r)
	return err
}

// userGetEndpoint returns a function to get the user with a given id
func userGetEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if err := validation.ValidateUuid(id); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := s.GetUser(context.Background(), id)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			e.HandleHttp(res, e.ErrNotFound.Error(), http.StatusNotFound)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		removeHashPass(user)
		json.NewEncoder(res).Encode(user)
	}
}

// userGetEmailEndpoint returns a function to get the user with a given email
func userGetEmailEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		email := chi.URLParam(req, "email")
		if err := validation.ValidateEmail(email); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := s.GetByEmail(context.Background(), email)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			e.HandleHttp(res, e.ErrNotFound.Error(), http.StatusNotFound)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		removeHashPass(user)
		json.NewEncoder(res).Encode(user)
	}
}

type userPostRequest struct {
	Email    string `json:"Email"`
	FullName string `json:"FullName"`
	Password string `json:"Password"`
}

// validate the userPostRequest
func (r *userPostRequest) validate() error {
	if r.FullName == "" || len(r.FullName) > 100 || len(r.FullName) < 3 {
		return e.ErrValidation
	}
	if r.Email == "" || len(r.Email) > 100 || len(r.Email) < 3 {
		return e.ErrValidation
	}
	if r.Password == "" || len(r.Password) > 100 || len(r.Password) < 3 {
		return e.ErrValidation
	}
	return nil
}

// userPostEndpoint returns a function to create a new user
func userPostEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userReq := userPostRequest{}
		err := json.NewDecoder(req.Body).Decode(&userReq)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = userReq.validate(); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}

		hashPass, err := hash.HashPassword(userReq.Password)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}

		user := models.User{
			Email:    userReq.Email,
			FullName: userReq.FullName,
			HashPass: hashPass,
		}
		err = s.SaveUser(context.Background(), &user)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		removeHashPass(&user)
		json.NewEncoder(res).Encode(&user)
	}
}

type userDeleteResponse struct {
	Id string `json:"Id"`
}

// userDeleteEndpoint returns a function to delete a user
func userDeleteEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if err := validation.ValidateUuid(id); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		userId, err := s.DeleteUser(context.Background(), id)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(userDeleteResponse{Id: userId})
	}
}

type userPutRequest struct {
	Email    string `json:"Email"`
	FullName string `json:"FullName"`
}

// validate the userPutRequest
func (r *userPutRequest) validate() error {
	if r.FullName == "" || len(r.FullName) > 100 || len(r.FullName) < 3 {
		return e.ErrValidation
	}
	if r.Email == "" || len(r.Email) > 100 || len(r.Email) < 3 {
		return e.ErrValidation
	}
	return nil
}

// userPutEndpoint returns a function to update a user
func userPutEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if err := validation.ValidateUuid(id); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		userReq := userPutRequest{}
		err := json.NewDecoder(req.Body).Decode(&userReq)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = userReq.validate(); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		user := models.User{
			Id:       id,
			Email:    userReq.Email,
			FullName: userReq.FullName,
		}
		err = s.UpdateUser(context.Background(), &user)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		removeHashPass(&user)
		json.NewEncoder(res).Encode(&user)
	}
}

type loginPostRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type loginPostResponse struct {
	Token string `json:"Token"`
}

// validate the loginPostRequest
func (req *loginPostRequest) validate() error {
	if req.Email == "" || len(req.Email) > 100 || len(req.Email) < 3 {
		return e.ErrValidation
	}
	if req.Password == "" || len(req.Password) > 100 || len(req.Password) < 3 {
		return e.ErrValidation
	}
	return nil
}

// loginPostEndpoint returns a function to login a user
func loginUserEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		loginReq := loginPostRequest{}
		err := json.NewDecoder(req.Body).Decode(&loginReq)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = loginReq.validate(); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := s.LoginUser(context.Background(), loginReq.Email, loginReq.Password)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(loginPostResponse{Token: token})
	}
}

// loginAdminEndpoint returns a function to login an admin
func loginAdminEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		loginReq := loginPostRequest{}
		err := json.NewDecoder(req.Body).Decode(&loginReq)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = loginReq.validate(); err != nil {
			e.HandleHttp(res, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := s.LoginAdmin(context.Background(), loginReq.Email, loginReq.Password)
		if err != nil {
			e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(loginPostResponse{Token: token})
	}
}

func removeHashPass(user *models.User) {
	user.HashPass = ""
}

type middlewareFunc func(http.Handler) http.Handler

// AdminOnly returns a function to check if the user is an admin
func AdminOnly(s service.Users) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("Authorization")
			log.Println("AdminOnly route !", token)
			isAuthorized, err := s.IsAdmin(req.Context(), token)
			if err != nil {
				e.HandleHttp(res, err.Error(), http.StatusInternalServerError)
				return
			}
			if !isAuthorized {
				e.HandleHttp(res, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(res, req)
		})
	}
}
