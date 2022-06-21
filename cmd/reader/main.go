package main

import (
	"time"

	"github.com/pseudoincorrect/bariot/internal/reader/service"
	"github.com/pseudoincorrect/bariot/internal/reader/ws"
	authClient "github.com/pseudoincorrect/bariot/pkg/auth/client"
	natsClient "github.com/pseudoincorrect/bariot/pkg/nats/client"
	thingsClient "github.com/pseudoincorrect/bariot/pkg/things/client"
	"github.com/pseudoincorrect/bariot/pkg/utils/env"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
)

func main() {
	conf := loadConfig()
	logger.Info("Reader service online")
	reader := createService()
	wsConfig := ws.Config{
		Host:    conf.readerWsHost,
		Port:    conf.readerWsPort,
		Service: reader,
	}
	ws.Start(wsConfig)
	for {
		time.Sleep(time.Millisecond * 100)
	}
}

type config struct {
	bariotEnv      string
	authGrpcHost   string
	authGrpcPort   string
	thingsGrpcHost string
	thingsGrpcPort string
	natsHost       string
	natsPort       string
	readerWsPort   string
	readerWsHost   string
}

// Load config from environment variables
func loadConfig() config {
	var conf = config{
		bariotEnv:      env.GetEnv("BARIOT_ENV"),
		authGrpcHost:   env.GetEnv("AUTH_GRPC_HOST"),
		authGrpcPort:   env.GetEnv("AUTH_GRPC_PORT"),
		thingsGrpcHost: env.GetEnv("THINGS_GRPC_HOST"),
		thingsGrpcPort: env.GetEnv("THINGS_GRPC_PORT"),
		natsHost:       env.GetEnv("NATS_HOST"),
		natsPort:       env.GetEnv("NATS_PORT"),
		readerWsHost:   env.GetEnv("READER_WS_HOST"),
		readerWsPort:   env.GetEnv("READER_WS_PORT"),
	}
	return conf
}

// createService with necessary clients
func createService() service.Reader {
	conf := loadConfig()
	authConf := authClient.Conf{Host: conf.authGrpcHost, Port: conf.authGrpcPort}
	auth := authClient.New(authConf)
	auth.StartAuthClient()
	thingsConf := thingsClient.Conf{Host: conf.thingsGrpcHost, Port: conf.thingsGrpcPort}
	things := thingsClient.New(thingsConf)
	things.StartThingsClient()
	natsConf := natsClient.Conf{Host: conf.natsHost, Port: conf.natsPort}
	nats := natsClient.New(natsConf)
	nats.Connect(natsClient.NatsSetupConnOptions("reader"))
	reader := service.New(&auth, &things, &nats)
	return &reader
}
