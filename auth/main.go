package main

import (
	"log"

	"github.com/pseudoincorrect/bariot/auth/rpc/server"
	"github.com/pseudoincorrect/bariot/auth/service"
	util "github.com/pseudoincorrect/bariot/auth/utilities"
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
		environment: util.GetEnv("BARIOT_ENV"),
		rpcHost:     util.GetEnv("RPC_HOST"),
		rpcPort:     util.GetEnv("RPC_PORT"),
		adminSecret: util.GetEnv("ADMIN_SECRET"),
		jwtSecret:   util.GetEnv("JWT_SECRET"),
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

	err := server.Start(serverConf)
	if err != nil {
		log.Panic("Error starting GRPC server:", err)
	}
}
