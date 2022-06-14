package mockReader

import (
	"github.com/pseudoincorrect/bariot/reader/service"
	"github.com/stretchr/testify/mock"
)

type MockReader struct {
	mock.Mock
}

var _ service.ReaderSvc = (*MockReader)(nil)

func (m *MockReader) AuthorizeSingleThing(userToken string, thingId string) (bool, error) {
	args := m.Called()
	return args.Bool(1), args.Error(0)
}
