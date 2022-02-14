package service

import (
	"context"
	"fmt"

	"github.com/pseudoincorrect/bariot/users/models"
)

type Users interface {
	SaveUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUser(ctx context.Context, id string) (*models.User, error)
	DeleteUser(ctx context.Context, id string) (string, error)
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
}

// type check on userService
var _ Users = (*usersService)(nil)

type usersService struct {
	repository models.UsersRepository
}

/// New creates a new user service
func New(repository models.UsersRepository) Users {
	return &usersService{repository}
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
