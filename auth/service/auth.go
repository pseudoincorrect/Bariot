package service

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	admin = "admin"
	user  = "user"
	thing = "thing"
)

type Auth interface {
	GetAdminToken() (string, error)
	GetUserToken(string) (string, error)
	GetThingToken(string) (string, error)
	ValidateToken(string) (bool, error)
	GetClaimsToken(string) (*AuthClaim, error)
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
	Role string `json:"Role"`
	jwt.StandardClaims
}

func (s *authService) GetAdminToken() (string, error) {
	token, err := s.makeToken(admin, "0", 1)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) GetUserToken(userId string) (string, error) {
	token, err := s.makeToken(user, userId, 24)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) GetThingToken(thingId string) (string, error) {
	token, err := s.makeToken(thing, thingId, 24)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *authService) makeToken(role string, subjet string, hours time.Duration) (string, error) {
	claims := AuthClaim{
		role,
		jwt.StandardClaims{
			Subject:   subjet,
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * hours).Unix(),
			Issuer:    s.environment,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secret)
	if err != nil {
		return "", err
	}
	log.Println(tokenString, err)
	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (bool, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaim{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return false, err
	}
	if _, ok := token.Claims.(*AuthClaim); ok && token.Valid {
		return true, nil
	}
	return false, jwt.ErrInvalidKey
}

func (s *authService) GetClaimsToken(tokenString string) (*AuthClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AuthClaim{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AuthClaim); ok && token.Valid {
		log.Println(claims)
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}
