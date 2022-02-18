package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Auth interface {
	GetAdminToken() (string, error)
	ValidateToken(string) (bool, error)
}

var _ Auth = (*authService)(nil)

type authService struct {
	secret      []byte
	environment string
}

type ServiceConf struct {
	Secret      string
	Environment string
}

func New(c ServiceConf) Auth {
	return &authService{
		secret:      []byte(c.Secret),
		environment: c.Environment,
	}
}

type AuthClaim struct {
	Role string `json:"foo"`
	jwt.StandardClaims
}

func (s *authService) GetAdminToken() (string, error) {
	claims := AuthClaim{
		"admin",
		jwt.StandardClaims{
			Subject:   "0",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    s.environment,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}
	fmt.Println(tokenString, err)
	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaim{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(*AuthClaim); ok && token.Valid {
		fmt.Println(claims.Role)
		return true, nil
	}
	return false, nil
}
