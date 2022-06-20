package client

import (
	"context"
	"testing"

	e "github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/tests/mocks/grpc"
	"github.com/stretchr/testify/assert"
)

const host = "localhost"
const port = "50051"

var cli thingsClient
var clientMock grpc.MockThings

func startClient() {
	conf := Conf{Host: host, Port: port}
	cli = New(conf)
	cli.StartThingsClient()
	clientMock = grpc.NewMockThings()
	cli.Client = &clientMock
}

func TestMain(m *testing.M) {
	startClient()
	m.Run()
}

func TestGetUserOfThing(t *testing.T) {
	ctx := context.Background()
	thingId := "000.000.001"
	userId := "000.000.002"
	clientMock.On("GetUserOfThing").Return(nil, userId).Once()
	res, err := cli.GetUserOfThing(ctx, thingId)
	assert.Nil(t, err, "should not throw an error")
	assert.Equal(t, userId, res, "should be the same user")
	clientMock.On("GetUserOfThing").Return(e.ErrGrpc, "").Once()
	_, err = cli.GetUserOfThing(ctx, thingId)
	assert.NotNil(t, err, "should throw an error")
}
