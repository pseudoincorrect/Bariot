package serviceMock

import (
	"context"

	"github.com/pseudoincorrect/bariot/users/models"
	"github.com/pseudoincorrect/bariot/users/service"
	"github.com/stretchr/testify/mock"
)

// Static type checking
var _ service.Users = (*Mock)(nil)

type Mock struct {
	mock.Mock
}

// SaveUser saves a user to repository with user model
func (m *Mock) SaveUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(user)

	gottenUser := makeUser(args.String(1), args.String(2))
	return &gottenUser, args.Error(0)
}

// GetUser returns a user from repository by id
func (m *Mock) GetUser(ctx context.Context, id string) (*models.User, error) {
	args := m.Called(id)
	if args.String(1) == "" {
		return nil, args.Error(0)
	}
	gottenUser := makeUser(args.String(1), args.String(2))
	return &gottenUser, args.Error(0)
}

// GetByEmail returns a user from repository by email
func (m *Mock) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(email)
	gottenUser := makeUser(args.String(1), args.String(2))
	return &gottenUser, args.Error(0)
}

// DeleteUser deletes a user from repository by id
func (m *Mock) DeleteUser(ctx context.Context, id string) (string, error) {
	args := m.Called(id)
	return args.String(1), args.Error(0)
}

// UpdateUser updates a user in repository by user model
func (m *Mock) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	args := m.Called(user)
	gottenUser := makeUser(args.String(1), args.String(2))
	return &gottenUser, args.Error(0)
}

func (m *Mock) LoginUser(ctx context.Context, email string, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(1), args.Error(0)
}

func (m *Mock) LoginAdmin(ctx context.Context, email string, password string) (string, error) {
	args := m.Called(email, password)
	return args.String(1), args.Error(0)
}

func (m *Mock) IsAdmin(ctx context.Context, token string) (bool, error) {
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
