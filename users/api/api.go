package api

import (
	"context"
	"encoding/json"
	"fmt"
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

func createEndpoint(s service.Users, r *chi.Mux) {
	// only admins can manage users
	adminGroup := r.Group(nil)
	adminGroup.Use(AdminOnly(s))
	adminGroup.Get("/{id}", userGetEndpoint(s))
	adminGroup.Post("/", userPostEndpoint(s))
	adminGroup.Delete("/{id}", userDeleteEndpoint(s))
	adminGroup.Put("/{id}", userPutEndpoint(s))
	r.Post("/login", loginUserEndpoint(s))
	r.Post("/login/admin", loginAdminEndpoint(s))
}

func startRouter(port string, r *chi.Mux) {
	addr := ":" + port
	http.ListenAndServe(addr, r)
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
		removeHashPass(user)
		json.NewEncoder(res).Encode(user)
	}
}

type userPostRequest struct {
	Email    string `json:"Email"`
	FullName string `json:"FullName"`
	Password string `json:"Password"`
}

func (r *userPostRequest) validate() error {
	if r.FullName == "" || len(r.FullName) > 100 || len(r.FullName) < 3 {
		return errors.NewValidationError("Invalid user name")
	}
	if r.Email == "" || len(r.Email) > 100 || len(r.Email) < 3 {
		return errors.NewValidationError("Invalid email")
	}
	if r.Password == "" || len(r.Password) > 100 || len(r.Password) < 3 {
		return errors.NewValidationError("Invalid Password")
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

		hashPass, err := utils.HashPassword(userReq.Password)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		user := models.User{
			Email:    userReq.Email,
			FullName: userReq.FullName,
			HashPass: hashPass,
		}
		savedUser, err := s.SaveUser(context.Background(), &user)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set("Content-Type", "application/json")
		removeHashPass(savedUser)
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
		removeHashPass(updatedUser)
		json.NewEncoder(res).Encode(updatedUser)
	}
}

type loginPostRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type loginPostResponse struct {
	Token string `json:"Token"`
}

func (req *loginPostRequest) validate() error {
	if req.Email == "" || len(req.Email) > 100 || len(req.Email) < 3 {
		return errors.NewValidationError("Invalid user email")
	}
	if req.Password == "" || len(req.Password) > 100 || len(req.Password) < 3 {
		return errors.NewValidationError("Invalid user password")
	}
	return nil
}

func loginUserEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		loginReq := loginPostRequest{}
		err := json.NewDecoder(req.Body).Decode(&loginReq)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = loginReq.validate(); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := s.LoginUser(context.Background(), loginReq.Email, loginReq.Password)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(loginPostResponse{Token: token})
	}
}

func loginAdminEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		loginReq := loginPostRequest{}
		err := json.NewDecoder(req.Body).Decode(&loginReq)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = loginReq.validate(); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := s.LoginAdmin(context.Background(), loginReq.Email, loginReq.Password)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(res).Encode(loginPostResponse{Token: token})
	}
}

func removeHashPass(user *models.User) {
	user.HashPass = ""
}

type middlewareFunc func(http.Handler) http.Handler

func AdminOnly(s service.Users) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			token := req.Header.Get("Authorization")
			fmt.Println("AdminOnly route !", token)
			isAuthorized, err := s.IsAdmin(req.Context(), token)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			if !isAuthorized {
				http.Error(res, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(res, req)
		})
	}
}
