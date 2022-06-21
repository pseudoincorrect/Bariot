package clients

import (
	natsGo "github.com/nats-io/nats.go"
	natsClient "github.com/pseudoincorrect/bariot/pkg/nats/client"
	"github.com/stretchr/testify/mock"
)

type NatsClientMock struct {
	mock.Mock
}

var _ natsClient.Nats = (*NatsClientMock)(nil)

func NewNatsClientMock() NatsClientMock {
	return NatsClientMock{}
}

func (m *NatsClientMock) Connect(opts []natsGo.Option) error {
	args := m.Called()
	return args.Error(0)
}

func (m *NatsClientMock) Disconnect() {
}

func (m *NatsClientMock) Subscribe(subject string, handler natsGo.MsgHandler) (*natsGo.Subscription, error) {
	args := m.Called()
	return nil, args.Error(0)
}

func (m *NatsClientMock) CreatePublisher() natsClient.NatsPubType {
	return nil
}

func (m *NatsClientMock) Publish(subject string, payload string) error {
	args := m.Called()
	return args.Error(0)
}
