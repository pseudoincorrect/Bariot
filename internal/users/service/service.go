package service

import (
	"context"
	"log"

	"github.com/pseudoincorrect/bariot/internal/users/hash"
	"github.com/pseudoincorrect/bariot/internal/users/models"
	authClient "github.com/pseudoincorrect/bariot/pkg/auth/client"
	"github.com/pseudoincorrect/bariot/pkg/errors"
)

type Users interface {
	SaveUser(ctx context.Context, model *models.User) (err error)
	GetUser(ctx context.Context, userId string) (model *models.User, err error)
	GetByEmail(ctx context.Context, email string) (model *models.User, err error)
	DeleteUser(ctx context.Context, userId string) (id string, err error)
	UpdateUser(ctx context.Context, model *models.User) (err error)
	LoginUser(ctx context.Context, email string, passW string) (token string, err error)
	LoginAdmin(ctx context.Context, email string, passW string) (token string, err error)
	IsAdmin(ctx context.Context, token string) (isAdmin bool, err error)
}

// Static type checking
var _ Users = (*usersService)(nil)

type usersService struct {
	repository models.UsersRepository
	auth       authClient.Auth
}

// New creates a new user service
func New(repository models.UsersRepository, auth authClient.Auth) Users {
	return &usersService{repository, auth}
}

// SaveUser saves a user to repository with user model
func (s *usersService) SaveUser(ctx context.Context, user *models.User) error {
	err := s.repository.Save(ctx, user)
	if err != nil {
		log.Println("Save User error:", err)
		return err
	}
	return nil
}

// GetUser returns a user from repository by id
func (s *usersService) GetUser(ctx context.Context, id string) (*models.User, error) {
	user, err := s.repository.Get(ctx, id)
	if err != nil {
		log.Println("Get User error:", err)
		return nil, err
	}
	return user, nil
}

// GetByEmail returns a user from repository by email
func (s *usersService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		log.Println("Get User by email error:", err)
		return nil, err
	}
	return user, nil
}

// DeleteUser deletes a user from repository by id
func (s *usersService) DeleteUser(ctx context.Context, id string) (string, error) {
	resId, err := s.repository.Delete(ctx, id)
	if err != nil {
		log.Println("Delete User error:", err)
		return "", err
	}
	return resId, nil
}

// UpdateUser updates a user in repository by user model
func (s *usersService) UpdateUser(ctx context.Context, user *models.User) error {
	err := s.repository.Update(ctx, user)
	if err != nil {
		log.Println("Update User error:", err)
		return err
	}
	return nil
}

func (s *usersService) LoginUser(ctx context.Context, email string, password string) (string, error) {
	user, err := s.GetByEmail(context.Background(), email)
	if err != nil {
		return "", errors.ErrDb
	}
	if user == nil {
		return "", errors.ErrNotFound
	}
	if !hash.CheckPasswordHash(password, user.HashPass) {
		return "", errors.ErrPassword
	}
	token, err := s.auth.GetUserToken(ctx, user.Id)
	if err != nil {
		return "", errors.ErrAuthn
	}
	return token, nil
}

func (s *usersService) LoginAdmin(ctx context.Context, email string, password string) (string, error) {
	user, err := s.GetByEmail(context.Background(), email)
	if err != nil {
		return "", errors.ErrDb
	}
	if user == nil {
		return "", errors.ErrNotFound
	}
	if !hash.CheckPasswordHash(password, user.HashPass) {
		return "", errors.ErrPassword
	}
	token, err := s.auth.GetAdminToken(ctx)
	if err != nil {
		return "", errors.ErrAuthn
	}
	return token, nil
}

func (s *usersService) IsAdmin(ctx context.Context, token string) (bool, error) {
	return s.auth.IsAdmin(ctx, token)
}
