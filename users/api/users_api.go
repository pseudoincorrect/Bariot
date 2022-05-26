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
	"github.com/pseudoincorrect/bariot/users/models"
	"github.com/pseudoincorrect/bariot/users/service"
	"github.com/pseudoincorrect/bariot/users/utilities/hash"
)

func InitApi(port string, s service.Users) error {
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

func startRouter(port string, r *chi.Mux) error {
	addr := ":" + port
	err := http.ListenAndServe(addr, r)
	return err
}

func userGetEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if err := validation.ValidateUuid(id); err != nil {
			appErr.Http(res, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := s.GetUser(context.Background(), id)
		if err != nil {
			appErr.Http(res, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			appErr.Http(res, appErr.ErrUserNotFound.Error(), http.StatusNotFound)
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
		return appErr.ErrValidation
	}
	if r.Email == "" || len(r.Email) > 100 || len(r.Email) < 3 {
		return appErr.ErrValidation
	}
	if r.Password == "" || len(r.Password) > 100 || len(r.Password) < 3 {
		return appErr.ErrValidation
	}
	return nil
}

func userPostEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userReq := userPostRequest{}
		err := json.NewDecoder(req.Body).Decode(&userReq)
		if err != nil {
			appErr.Http(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = userReq.validate(); err != nil {
			appErr.Http(res, err.Error(), http.StatusBadRequest)
			return
		}

		hashPass, err := hash.HashPassword(userReq.Password)
		if err != nil {
			appErr.Http(res, err.Error(), http.StatusInternalServerError)
			return
		}

		user := models.User{
			Email:    userReq.Email,
			FullName: userReq.FullName,
			HashPass: hashPass,
		}
		savedUser, err := s.SaveUser(context.Background(), &user)
		if err != nil {
			appErr.Http(res, err.Error(), http.StatusInternalServerError)
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
		if err := validation.ValidateUuid(id); err != nil {
			appErr.Http(res, err.Error(), http.StatusBadRequest)
			return
		}
		userId, err := s.DeleteUser(context.Background(), id)
		if err != nil {
			appErr.Http(res, err.Error(), http.StatusInternalServerError)
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
		return appErr.ErrValidation
	}
	if r.Email == "" || len(r.Email) > 100 || len(r.Email) < 3 {
		return appErr.ErrValidation
	}
	return nil
}

func userPutEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := chi.URLParam(req, "id")
		if err := validation.ValidateUuid(id); err != nil {
			appErr.Http(res, err.Error(), http.StatusBadRequest)
			return
		}
		userReq := userPutRequest{}
		err := json.NewDecoder(req.Body).Decode(&userReq)
		if err != nil {
			appErr.Http(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = userReq.validate(); err != nil {
			appErr.Http(res, err.Error(), http.StatusBadRequest)
			return
		}
		user := models.User{
			Id:       id,
			Email:    userReq.Email,
			FullName: userReq.FullName,
		}
		updatedUser, err := s.UpdateUser(context.Background(), &user)
		if err != nil {
			appErr.Http(res, err.Error(), http.StatusInternalServerError)
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
		return appErr.ErrValidation
	}
	if req.Password == "" || len(req.Password) > 100 || len(req.Password) < 3 {
		return appErr.ErrValidation
	}
	return nil
}

func loginUserEndpoint(s service.Users) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		loginReq := loginPostRequest{}
		err := json.NewDecoder(req.Body).Decode(&loginReq)
		if err != nil {
			appErr.Http(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = loginReq.validate(); err != nil {
			appErr.Http(res, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := s.LoginUser(context.Background(), loginReq.Email, loginReq.Password)
		if err != nil {
			appErr.Http(res, err.Error(), http.StatusInternalServerError)
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
			appErr.Http(res, err.Error(), http.StatusBadRequest)
			return
		}
		if err = loginReq.validate(); err != nil {
			appErr.Http(res, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := s.LoginAdmin(context.Background(), loginReq.Email, loginReq.Password)
		if err != nil {
			appErr.Http(res, err.Error(), http.StatusInternalServerError)
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
			log.Println("AdminOnly route !", token)
			isAuthorized, err := s.IsAdmin(req.Context(), token)
			if err != nil {
				appErr.Http(res, err.Error(), http.StatusInternalServerError)
				return
			}
			if !isAuthorized {
				appErr.Http(res, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(res, req)
		})
	}
}
