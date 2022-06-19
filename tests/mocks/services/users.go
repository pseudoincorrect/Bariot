package services

import (
	"context"

	"github.com/pseudoincorrect/bariot/internal/users/models"
	"github.com/pseudoincorrect/bariot/internal/users/service"
	"github.com/stretchr/testify/mock"
)

// Static type checking
var _ service.Users = (*UsersMock)(nil)

type UsersMock struct {
	mock.Mock
}

func NewUsersMock() UsersMock {
	return UsersMock{}
}

// SaveUser saves a user to repository with user model
func (m *UsersMock) SaveUser(ctx context.Context, user *models.User) error {
	args := m.Called(user)

	*user = makeUser(args.String(1), args.String(2))
	return args.Error(0)
}

// GetUser returns a user from repository by id
func (m *UsersMock) GetUser(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(id)
	if args.String(1) == "" {
		return nil, args.Error(0)
	}
	gottenUser := makeUser(args.String(1), args.String(2))
	return &gottenUser, args.Error(0)
}

// GetByEmail returns a user from repository by email
func (m *UsersMock) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(email)
	gottenUser := makeUser(args.String(1), args.String(2))
	return &gottenUser, args.Error(0)
}

// DeleteUser deletes a user from repository by id
func (m *UsersMock) DeleteUser(ctx context.Context, id string) (string, error) {
	args := m.Called(id)
	return args.String(1), args.Error(0)
}

// UpdateUser updates a user in repository by user model
func (m *UsersMock) UpdateUser(ctx context.Context, user *models.User) error {
	args := m.Called(user)
	*user = makeUser(args.String(1), args.String(2))
	return args.Error(0)
}

func (m *UsersMock) LoginUser(ctx context.Context, email string, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(1), args.Error(0)
}

func (m *UsersMock) LoginAdmin(ctx context.Context, email string, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(1), args.Error(0)
}

func (m *UsersMock) IsAdmin(ctx context.Context, token string) (bool, error) {
	args := m.Called(token)
	return args.Bool(1), args.Error(0)
}

func makeUser(id string, name string) models.User {
	return models.User{
		Id:        id,
		CreatedAt: "",
		Email:     name + "@test.com",
		FullName:  name,
		HashPass:  "ldzkjhgbksdqlkmfdsljkfsmlk",
		Metadata:  models.Metadata{"address": "Cornimont"},
	}
}
