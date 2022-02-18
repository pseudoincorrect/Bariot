package server

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/pseudoincorrect/bariot/auth/rpc/auth"
	"google.golang.org/grpc"

	"github.com/pseudoincorrect/bariot/auth/service"
)

type server struct {
	pb.UnimplementedAuthServer
	AuthService service.Auth
}

func (s *server) GetAdminToken(ctx context.Context, in *pb.GetAdminTokenRequest) (*pb.GetAdminTokenResponse, error) {
	token, err := s.AuthService.GetAdminToken()
	if err != nil {
		return nil, err
	}
	return &pb.GetAdminTokenResponse{Jwt: token}, nil
}

func (s *server) ValidateToken(ctx context.Context, in *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	valid, err := s.AuthService.ValidateToken(in.Jwt)
	if err != nil {
		return nil, err
	}
	return &pb.ValidateTokenResponse{Valid: valid}, nil
}

type ServerConf struct {
	AuthService service.Auth
	Port        string
}

func Start(c ServerConf) error {
	addr := ":" + c.Port
	fmt.Println("Starting Auth GRPC on", addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}
	s := grpc.NewServer()
	pb.RegisterAuthServer(s, &server{pb.UnimplementedAuthServer{}, c.AuthService})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
		return err
	}
	return nil
}
