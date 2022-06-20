package grpc

import (
	"context"

	pb "github.com/pseudoincorrect/bariot/pkg/things/grpc"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockThings struct {
	mock.Mock
}

var _ pb.ThingsClient = (*MockThings)(nil)

func NewMockThings() MockThings {
	return MockThings{}
}

func (m *MockThings) GetUserOfThing(ctx context.Context, in *pb.GetUserOfThingRequest, opts ...grpc.CallOption) (*pb.GetUserOfThingResponse, error) {
	args := m.Called()
	userId := args.String(1)
	res := pb.GetUserOfThingResponse{
		UserId: userId,
	}
	return &res, args.Error(0)
}
