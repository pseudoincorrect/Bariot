package service

import (
	"context"
	"fmt"

	"github.com/pseudoincorrect/bariot/users/models"
	"github.com/pseudoincorrect/bariot/users/rpc/auth/client"
	utils "github.com/pseudoincorrect/bariot/users/utilities"
	"github.com/pseudoincorrect/bariot/users/utilities/errors"
)

type Users interface {
	SaveUser(context.Context, *models.User) (*models.User, error)
	GetUser(context.Context, string) (*models.User, error)
	GetByEmail(context.Context, string) (*models.User, error)
	DeleteUser(context.Context, string) (string, error)
	UpdateUser(context.Context, *models.User) (*models.User, error)
	Login(context.Context, string, string) (string, error)
}

// type check on userService
var _ Users = (*usersService)(nil)

type usersService struct {
	repository models.UsersRepository
	auth       client.Auth
}

/// New creates a new user service
func New(repository models.UsersRepository, auth client.Auth) Users {
	return &usersService{repository, auth}
}

/// SaveUser saves a user to repository with user model
func (s *usersService) SaveUser(ctx context.Context, user *models.User) (*models.User, error) {
	savedUser, err := s.repository.Save(ctx, user)
	if err != nil {
		fmt.Println("Save User error:", err)
		return nil, err
	}
	return savedUser, nil
}

/// GetUser returns a user from repository by id
func (s *usersService) GetUser(ctx context.Context, id string) (*models.User, error) {
	user, err := s.repository.Get(ctx, id)
	if err != nil {
		fmt.Println("Get User error:", err)
		return nil, err
	}
	return user, nil
}

func (s *usersService) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		fmt.Println("Get User by email error:", err)
		return nil, err
	}
	return user, nil
}

/// DeleteUser deletes a user from repository by id
func (s *usersService) DeleteUser(ctx context.Context, id string) (string, error) {
	resId, err := s.repository.Delete(ctx, id)
	if err != nil {
		fmt.Println("Delete User error:", err)
		return "", err
	}
	return resId, nil
}

/// UpdateUser updates a user in repository by user model
func (s *usersService) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	updatedUser, err := s.repository.Update(ctx, user)
	if err != nil {
		fmt.Println("Update User error:", err)
		return nil, err
	}
	return updatedUser, nil
}

func (s *usersService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.GetByEmail(context.Background(), email)
	if err != nil {
		return "", errors.NewDbError(err.Error())
	}
	if user == nil {
		return "", errors.NewUserNotFoundError(email)
	}
	if !utils.CheckPasswordHash(password, user.HashPass) {
		return "", errors.NewPasswordError()
	}
	token, err := s.auth.GetUserToken(ctx, user.Id)
	if err != nil {
		return "", errors.NewAuthError(err.Error())
	}
	return token, nil
}
