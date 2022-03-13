package service

import (
	"context"
	"log"

	authClient "github.com/pseudoincorrect/bariot/pkg/auth/client"
	"github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/users/models"
	"github.com/pseudoincorrect/bariot/users/utilities/hash"
)

type ctxt context.Context

type Users interface {
	SaveUser(ctxt, *models.User) (*models.User, error)
	GetUser(ctxt, string) (*models.User, error)
	GetByEmail(ctxt, string) (*models.User, error)
	DeleteUser(ctxt, string) (string, error)
	UpdateUser(ctxt, *models.User) (*models.User, error)
	LoginUser(ctxt, string, string) (string, error)
	LoginAdmin(ctxt, string, string) (string, error)
	IsAdmin(ctxt, string) (bool, error)
}

// type check on userService
var _ Users = (*usersService)(nil)

type usersService struct {
	repository models.UsersRepository
	auth       authClient.Auth
}

/// New creates a new user service
func New(repository models.UsersRepository, auth authClient.Auth) Users {
	return &usersService{repository, auth}
}

/// SaveUser saves a user to repository with user model
func (s *usersService) SaveUser(ctx ctxt, user *models.User) (*models.User, error) {
	savedUser, err := s.repository.Save(ctx, user)
	if err != nil {
		log.Println("Save User error:", err)
		return nil, err
	}
	return savedUser, nil
}

/// GetUser returns a user from repository by id
func (s *usersService) GetUser(ctx ctxt, id string) (*models.User, error) {
	user, err := s.repository.Get(ctx, id)
	if err != nil {
		log.Println("Get User error:", err)
		return nil, err
	}
	return user, nil
}

func (s *usersService) GetByEmail(ctx ctxt, email string) (*models.User, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		log.Println("Get User by email error:", err)
		return nil, err
	}
	return user, nil
}

/// DeleteUser deletes a user from repository by id
func (s *usersService) DeleteUser(ctx ctxt, id string) (string, error) {
	resId, err := s.repository.Delete(ctx, id)
	if err != nil {
		log.Println("Delete User error:", err)
		return "", err
	}
	return resId, nil
}

/// UpdateUser updates a user in repository by user model
func (s *usersService) UpdateUser(ctx ctxt, user *models.User) (*models.User, error) {
	updatedUser, err := s.repository.Update(ctx, user)
	if err != nil {
		log.Println("Update User error:", err)
		return nil, err
	}
	return updatedUser, nil
}

func (s *usersService) LoginUser(ctx ctxt, email string, password string) (string, error) {
	user, err := s.GetByEmail(context.Background(), email)
	if err != nil {
		return "", errors.ErrDb
	}
	if user == nil {
		return "", errors.ErrUserNotFound
	}
	if !hash.CheckPasswordHash(password, user.HashPass) {
		return "", errors.ErrPassword
	}
	token, err := s.auth.GetUserToken(ctx, user.Id)
	if err != nil {
		return "", errors.ErrAuthentication
	}
	return token, nil
}

func (s *usersService) LoginAdmin(ctx ctxt, email string, password string) (string, error) {
	user, err := s.GetByEmail(context.Background(), email)
	if err != nil {
		return "", errors.ErrDb
	}
	if user == nil {
		return "", errors.ErrUserNotFound
	}
	if !hash.CheckPasswordHash(password, user.HashPass) {
		return "", errors.ErrPassword
	}
	token, err := s.auth.GetAdminToken(ctx)
	if err != nil {
		return "", errors.ErrAuthentication
	}
	return token, nil
}

func (s *usersService) IsAdmin(ctx ctxt, token string) (bool, error) {
	return s.auth.IsAdmin(ctx, token)
}
