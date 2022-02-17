package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/pseudoincorrect/bariot/auth/rpc/auth"
	util "github.com/pseudoincorrect/bariot/auth/utilities"
	"google.golang.org/grpc"
)

type config struct {
	rpcHost     string
	rpcPort     string
	adminSecret string
	jwtSecret   string
}

func loadConfig() config {
	var conf = config{
		rpcHost:     util.GetEnv("RPC_HOST"),
		rpcPort:     util.GetEnv("RPC_PORT"),
		adminSecret: util.GetEnv("ADMIN_SECRET"),
		jwtSecret:   util.GetEnv("JWT_SECRET"),
	}
	return conf
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedAuthServer
}

func (s *server) GetAdminToken(ctx context.Context, in *pb.GetAdminTokenRequest) (*pb.GetAdminTokenResponse, error) {
	fmt.Println("Got a GetAdminToken Request")
	return &pb.GetAdminTokenResponse{Jwt: "admin_token"}, nil
}

func main() {
	fmt.Println("Auth service...")
	conf := loadConfig()

	flag.Parse()
	// addr := conf.rpcHost + ":" + conf.rpcPort
	addr := ":" + conf.rpcPort

	fmt.Println("Starting Auth GRPC on", addr)

	lis, err := net.Listen("tcp", addr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterAuthServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
