package clients

import (
	"context"

	"github.com/pseudoincorrect/bariot/pkg/things/client"
	"github.com/stretchr/testify/mock"
)

type ThingsClientMock struct {
	mock.Mock
}

var _ client.Things = (*ThingsClientMock)(nil)

func NewThingsClientMock() ThingsClientMock {
	return ThingsClientMock{}
}

func (m *ThingsClientMock) StartThingsClient() error {
	args := m.Called()
	return args.Error(0)
}

func (m *ThingsClientMock) GetUserOfThing(ctx context.Context, thingId string) (userId string, err error) {
	args := m.Called()
	return args.String(1), args.Error(0)
}
