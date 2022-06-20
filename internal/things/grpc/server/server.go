package server

import (
	"context"
	"log"
	"net"

	tGrpc "github.com/pseudoincorrect/bariot/pkg/things/grpc"
	"github.com/pseudoincorrect/bariot/pkg/utils/debug"
	"google.golang.org/grpc"

	"github.com/pseudoincorrect/bariot/internal/things/service"
)

type server struct {
	Service service.Things
	tGrpc.UnimplementedThingsServer
}

// GetAdminToken returns a JWT token for the admin user
func (s *server) GetUserOfThing(ctx context.Context,
	in *tGrpc.GetUserOfThingRequest) (*tGrpc.GetUserOfThingResponse, error) {
	userId, err := s.Service.GetUserOfThing(ctx, in.ThingId)
	if err != nil {
		return nil, err
	}
	return &tGrpc.GetUserOfThingResponse{UserId: userId}, nil
}

type ServerConf struct {
	Service service.Things
	Port    string
}

// Start starts the GRPC server
func Start(conf ServerConf) (*grpc.Server, error) {
	addr := ":" + conf.Port
	debug.LogInfo("Starting Auth GRPC on", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return nil, err
	}
	s := grpc.NewServer()
	routes := server{conf.Service, tGrpc.UnimplementedThingsServer{}}
	tGrpc.RegisterThingsServer(s, &routes)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return nil, err
	}
	return s, nil
}
