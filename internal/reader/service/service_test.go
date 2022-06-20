package service

import (
	"testing"

	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/tests/mocks/clients"
	"github.com/stretchr/testify/assert"
)

var authCli clients.AuthClientMock
var thingsCli clients.ThingsClientMock
var natsCli clients.NatsClientMock

func newService() Reader {
	authCli = clients.NewAuthClientMock()
	thingsCli = clients.NewThingsClientMock()
	natsCli = clients.NewNatsClientMock()
	reader := New(&authCli, &thingsCli, &natsCli)
	return &reader
}

func TestAuthorizeSingleThing(tt *testing.T) {
	thingId := "000.000.001"
	userId := "000.000.002"
	wrongUserId := "000.000.003"
	s := newService()

	tt.Run("all goes well", func(t *testing.T) {
		authCli.UserId = userId
		thingsCli.On("GetUserOfThing").Return(nil, userId).Once()
		err := s.AuthorizeSingleThing("token", thingId)
		assert.Nil(t, err, "should not return an error")
	})

	tt.Run("wrong user id", func(t *testing.T) {
		authCli.UserId = wrongUserId
		thingsCli.On("GetUserOfThing").Return(nil, userId).Once()
		err := s.AuthorizeSingleThing("token", thingId)
		assert.Equal(t, err, e.ErrAuthz, "should return unauthorized error")
	})
}

// func TestReceiveThingData(tt *testing.T) {
// 	s := newService()
// 	thingId := "000.000.001"
// 	stop := make(chan bool)
// 	handler := func(msg string) {
// 	}

// 	tt.Run("all goes well", func(t *testing.T) {
// 		err := s.ReceiveThingData(thingId, handler, stop)
// 		assert.Nil(t, err, "should not throw an error")
// 	})
// }
