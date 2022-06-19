package clients

import (
	"context"
	"errors"

	auth "github.com/pseudoincorrect/bariot/pkg/auth/client"
)

const (
	AdminRole = "admin"
	UserRole  = "user"
	ThingRole = "thing"
)

// Static type checking
var _ auth.Auth = (*AuthClientMock)(nil)

func NewAuthClientMock() AuthClientMock {
	return AuthClientMock{}
}

type AuthClientMock struct {
	ThrowErr string
	UserId   string
	UserRole string
}

// StartAuthClient starts the auth client GRPC server
func (c *AuthClientMock) StartAuthClient() error {
	if c.ThrowErr != "" {
		return errors.New(c.ThrowErr)
	}
	return nil
}

// IsAdmin checks if the user is an admin given a token
func (c *AuthClientMock) IsAdmin(ctx context.Context, jwt string) (bool, error) {
	if c.ThrowErr != "" {
		return false, errors.New(c.ThrowErr)
	}
	return false, nil
}

// IsWhichUser checks if the user is a user given a token return role, user id
func (c *AuthClientMock) IsWhichUser(ctx context.Context, jwt string) (string, string, error) {
	if c.ThrowErr != "" {
		return "", "", errors.New(c.ThrowErr)
	}
	return c.UserRole, c.UserId, nil
}

// IsWhichThing whom a thing belong to, given a token
func (c *AuthClientMock) IsWhichThing(ctx context.Context, jwt string) (string, error) {
	if c.ThrowErr != "" {
		return "", errors.New(c.ThrowErr)
	}
	return "", nil
}

// GetThingToken returns a token for a thing given a user id and a thing id
func (c *AuthClientMock) GetThingToken(ctx context.Context, thingId string, userId string) (string, error) {
	if c.ThrowErr != "" {
		return "", errors.New(c.ThrowErr)
	}
	return "", nil
}

// GetAdminToken returns a token for an admin
func (c *AuthClientMock) GetAdminToken(ctx context.Context) (string, error) {
	if c.ThrowErr != "" {
		return "", errors.New(c.ThrowErr)
	}
	return "", nil
}

// GetUserToken returns a token for a user
func (c *AuthClientMock) GetUserToken(ctx context.Context, userId string) (string, error) {
	if c.ThrowErr != "" {
		return "", errors.New(c.ThrowErr)
	}
	return "", nil
}
