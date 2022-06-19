package service

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

const secret string = "superSecret"
const otherSecret string = "worseSecret"
const env string = "test"

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestGetAdminToken(t *testing.T) {
	srvc := makeService(secret)
	token, err := srvc.GetAdminToken()
	assert.Nil(t, err, "should get token without error")
	claims, err := srvc.GetClaimsUserToken(token)
	assert.Nil(t, err, "should get claims without error")
	assert.Equal(t, claims.Role, AdminRole, "Should be an admin")
}

func TestGetUserToken(t *testing.T) {
	srvc := makeService(secret)
	userId := "000.000.001"
	token, err := srvc.GetUserToken(userId)
	assert.Nil(t, err, "should get token without error")
	claims, err := srvc.GetClaimsUserToken(token)
	assert.Nil(t, err, "should get claims without error")
	assert.Equal(t, UserRole, claims.Role, "Should be an admin")
	assert.Equal(t, userId, claims.Subject, "Should be an admin")
}

func TestGetThingToken(t *testing.T) {
	srvc := makeService(secret)
	thingId := "000.000.001"
	userId := "000.000.002"
	token, err := srvc.GetThingToken(thingId, userId)
	assert.Nil(t, err, "should get token without error")
	claims, err := srvc.GetClaimsThingToken(token)
	assert.Nil(t, err, "should get claims without error")
	assert.Equal(t, thingId, claims.Subject, "Should be an admin")
}

func TestValidateUserToken(t *testing.T) {
	srvc := makeService(secret)
	userId := "000.000.001"
	token, err := srvc.GetUserToken(userId)
	assert.Nil(t, err, "should get token without error")
	err = srvc.ValidateUserToken(token)
	assert.Nil(t, err, "should validate token without error")
	otherSrvc := makeService(otherSecret)
	err = otherSrvc.ValidateUserToken(token)
	assert.NotNil(t, err, "should validate token without error")
	assert.Equal(t, jwt.ErrSignatureInvalid, err, "token should not be valid")
}

func TestValidateThingToken(t *testing.T) {
	srvc := makeService(secret)
	thingId := "000.000.001"
	userId := "000.000.002"
	token, err := srvc.GetThingToken(thingId, userId)
	assert.Nil(t, err, "should get token without error")
	err = srvc.ValidateThingToken(token)
	assert.Nil(t, err, "should validate token without error")
	otherSrvc := makeService(otherSecret)
	err = otherSrvc.ValidateThingToken(token)
	assert.NotNil(t, err, "should validate token without error")
	assert.Equal(t, jwt.ErrSignatureInvalid, err, "token should not be valid")
}

func makeService(theSecret string) Auth {
	conf := ServiceConf{
		Secret:      theSecret,
		Environment: env,
	}
	return New(conf)
}
