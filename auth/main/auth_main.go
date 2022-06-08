package main

import (
	"log"

	"github.com/pseudoincorrect/bariot/auth/rpc/server"
	"github.com/pseudoincorrect/bariot/auth/service"
	"github.com/pseudoincorrect/bariot/pkg/env"
)

type config struct {
	environment string
	rpcHost     string
	rpcPort     string
	adminSecret string
	jwtSecret   string
}

func loadConfig() config {
	var conf = config{
		environment: env.GetEnv("BARIOT_ENV"),
		rpcHost:     env.GetEnv("RPC_HOST"),
		rpcPort:     env.GetEnv("RPC_PORT"),
		adminSecret: env.GetEnv("ADMIN_SECRET"),
		jwtSecret:   env.GetEnv("JWT_SECRET"),
	}
	return conf
}

func main() {
	log.Println("Auth service...")

	conf := loadConfig()

	serviceConf := service.ServiceConf{
		Secret:      conf.jwtSecret,
		Environment: conf.environment,
	}
	service := service.New(serviceConf)

	serverConf := server.ServerConf{
		AuthService: service,
		Port:        conf.rpcPort,
	}

	_, err := server.Start(serverConf)
	if err != nil {
		log.Panic("Error starting GRPC server:", err)
	}
}
