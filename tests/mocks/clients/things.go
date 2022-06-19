package clients

import (
	"context"

	"github.com/pseudoincorrect/bariot/pkg/things/client"
	"github.com/stretchr/testify/mock"
)

type MockThingsClient struct {
	mock.Mock
}

var _ client.Things = (*MockThingsClient)(nil)

func NewMockThingsClient() MockThingsClient {
	return MockThingsClient{}
}

func (m *MockThingsClient) StartThingsClient() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockThingsClient) GetUserOfThing(ctx context.Context, thingId string) (userId string, err error) {
	args := m.Called()
	return args.String(1), args.Error(0)
}
