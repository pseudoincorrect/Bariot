package main

import (
	"log"

	"github.com/pseudoincorrect/bariot/internal/auth/grpc/server"
	"github.com/pseudoincorrect/bariot/internal/auth/service"
	"github.com/pseudoincorrect/bariot/pkg/env"
	"github.com/pseudoincorrect/bariot/pkg/utils/debug"
)

func main() {
	debug.LogInfo("Auth service...")
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
		rpcHost:     env.GetEnv("AUTH_GRPC_HOST"),
		rpcPort:     env.GetEnv("AUTH_GRPC_PORT"),
		adminSecret: env.GetEnv("ADMIN_SECRET"),
		jwtSecret:   env.GetEnv("JWT_SECRET"),
	}
	return conf
}
