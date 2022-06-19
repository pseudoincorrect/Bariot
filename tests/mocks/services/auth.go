package services

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pseudoincorrect/bariot/internal/auth/service"
	"github.com/stretchr/testify/mock"
)

type AuthMock struct {
	mock.Mock
}

var _ service.Auth = (*AuthMock)(nil)

func NewAuthMock() AuthMock {
	return AuthMock{}
}

// GetAdminToken returns a token for the admin user
func (m *AuthMock) GetAdminToken() (string, error) {
	args := m.Called()
	return args.String(1), args.Error(0)
}

// GetUserToken returns a token for the regular user
func (m *AuthMock) GetUserToken(userId string) (string, error) {
	args := m.Called(userId)
	return args.String(1), args.Error(0)
}

//GetThingsToken returns a token for the thing
func (m *AuthMock) GetThingToken(thingId string, userId string) (string, error) {
	args := m.Called(thingId, userId)
	return args.String(1), args.Error(0)
}

// ValidateUserToken returns true if the token is valid
func (m *AuthMock) ValidateUserToken(tokenString string) error {
	args := m.Called(tokenString)
	return args.Error(0)
}

//ValidateThingToken returns true if the token is valid
func (m *AuthMock) ValidateThingToken(tokenString string) error {
	args := m.Called(tokenString)
	return args.Error(0)
}

// GetClaimsUserToken returns the claims of the user token
func (m *AuthMock) GetClaimsUserToken(
	tokenString string) (*service.UserAuthClaim, error) {
	args := m.Called(tokenString)
	claims := service.UserAuthClaim{
		Role: args.String(1),
		StandardClaims: jwt.StandardClaims{
			Subject:   args.String(2),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(args.Int(3))).Unix(),
			Issuer:    "Test",
		},
	}
	return &claims, args.Error(0)
}

// GetClaimsThingToken returns the claims of the thing token
func (m *AuthMock) GetClaimsThingToken(
	tokenString string) (*service.ThingAuthClaim, error) {
	args := m.Called(tokenString)
	claims := service.ThingAuthClaim{
		UserId: args.String(1),
		StandardClaims: jwt.StandardClaims{
			Subject:   args.String(2),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(args.Int(3))).Unix(),
			Issuer:    "Test",
		},
	}
	return &claims, args.Error(0)
}

// // makeUserToken create and return a token for the regular user
// func (m *AuthMock) makeUserToken(userId string, role string, hours time.Duration) (string, error) {
// 	return "", nil
// }

// // makeThingToken create and return a token for a thing
// func (m *AuthMock) makeThingToken(thingId string, userId string, hours time.Duration) (string, error) {
// 	return "", nil
// }
