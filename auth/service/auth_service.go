package service

import (
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	AdminRole = "admin"
	UserRole  = "user"
)

type Auth interface {
	GetAdminToken() (string, error)
	GetUserToken(userId string) (string, error)
	GetThingToken(thingId string, userId string) (string, error)
	GetClaimsUserToken(token string) (*UserAuthClaim, error)
	GetClaimsThingToken(token string) (*ThingAuthClaim, error)
	ValidateThingToken(token string) error
	ValidateUserToken(token string) error
}

// Static type checking
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

type UserAuthClaim struct {
	Role string `json:"Role"`
	jwt.StandardClaims
}

type ThingAuthClaim struct {
	UserId string `json:"UserId"`
	jwt.StandardClaims
}

// GetAdminToken returns a token for the admin user
func (s *authService) GetAdminToken() (string, error) {
	token, err := s.makeUserToken("0", AdminRole, 1)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetUserToken returns a token for the regular user
func (s *authService) GetUserToken(userId string) (string, error) {
	token, err := s.makeUserToken(userId, UserRole, 24)
	if err != nil {
		return "", err
	}
	return token, nil
}

//GetThingsToken returns a token for the thing
func (s *authService) GetThingToken(thingId string, userId string) (string, error) {
	token, err := s.makeThingToken(thingId, userId, 24)
	if err != nil {
		return "", err
	}
	return token, nil
}

// makeUserToken create and return a token for the regular user
func (s *authService) makeUserToken(userId string, role string, hours time.Duration) (string, error) {
	claims := UserAuthClaim{
		Role: role,
		StandardClaims: jwt.StandardClaims{
			Subject:   userId,
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
	// log.Println(tokenString, err)
	return tokenString, nil
}

// makeThingToken create and return a token for a thing
func (s *authService) makeThingToken(thingId string, userId string, hours time.Duration) (string, error) {
	claims := ThingAuthClaim{
		UserId: userId,
		StandardClaims: jwt.StandardClaims{
			Subject:   thingId,
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
	// log.Println(tokenString, err)
	return tokenString, nil
}

// ValidateUserToken returns true if the token is valid
func (s *authService) ValidateUserToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &UserAuthClaim{}, s.jwtSecretGetter())
	if err != nil {
		return jwt.ErrSignatureInvalid
	}
	if _, ok := token.Claims.(*UserAuthClaim); ok && token.Valid {
		return nil
	}
	return jwt.ErrInvalidKey
}

//ValidateThingToken returns true if the token is valid
func (s *authService) ValidateThingToken(tokenString string) error {
	token, err := jwt.ParseWithClaims(tokenString, &ThingAuthClaim{}, s.jwtSecretGetter())
	if err != nil {
		return jwt.ErrSignatureInvalid
	}
	if _, ok := token.Claims.(*ThingAuthClaim); ok && token.Valid {
		return nil
	}
	return jwt.ErrInvalidKey
}

// GetClaimsUserToken returns the claims of the user token
func (s *authService) GetClaimsUserToken(tokenString string) (*UserAuthClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserAuthClaim{}, s.jwtSecretGetter())
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*UserAuthClaim); ok && token.Valid {
		// log.Println(claims)
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}

// GetClaimsThingToken returns the claims of the thing token
func (s *authService) GetClaimsThingToken(tokenString string) (*ThingAuthClaim, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ThingAuthClaim{}, s.jwtSecretGetter())
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*ThingAuthClaim); ok && token.Valid {
		// log.Println(claims)
		return claims, nil
	}
	return nil, jwt.ErrInvalidKey
}

func (s *authService) jwtSecretGetter() jwt.Keyfunc {
	f := func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	}
	return f
}
