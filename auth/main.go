package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/pseudoincorrect/bariot/auth/rpc/proto/authProto"
	"google.golang.org/grpc"
)

type config struct {
	httpPort   string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string
	dbName     string
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

func main() {
	fmt.Println("Auth service...")

	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	authProto.RegisterRouteGuideServer(grpcServer, newServer())
	grpcServer.Serve(lis)
	proto.Equal()
}
