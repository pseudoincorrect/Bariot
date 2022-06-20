package main

import (
	"log"
	"time"

	"github.com/pseudoincorrect/bariot/internal/things/db"
	"github.com/pseudoincorrect/bariot/internal/things/grpc/server"
	"github.com/pseudoincorrect/bariot/internal/things/http"
	"github.com/pseudoincorrect/bariot/internal/things/service"
	authClient "github.com/pseudoincorrect/bariot/pkg/auth/client"
	cdb "github.com/pseudoincorrect/bariot/pkg/cache"
	"github.com/pseudoincorrect/bariot/pkg/env"
	e "github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/debug"
)

func main() {
	conf := loadConfig()
	debug.LogInfo("Things service online")
	thingsService, err := createService()
	if err != nil {
		log.Panic("Create service error:", err)
	}
	grpcServerConf := server.ServerConf{
		Service: thingsService,
		Port:    conf.thingsGrpcPort,
	}
	go func() {
		_, err = server.Start(grpcServerConf)
		if err != nil {
			err = e.Handle(e.ErrHttpServer, err, "starting http server")
			log.Panic(err)
		}
	}()
	go func() {
		err = startHttp(thingsService)
		if err != nil {
			err = e.Handle(e.ErrGrpcServer, err, "starting grpc server")
			log.Panic(err)
		}
	}()

	for {
		time.Sleep(time.Millisecond * 100)
	}
}

type config struct {
	httpPort       string
	authGrpcHost   string
	authGrpcPort   string
	thingsGrpcHost string
	thingsGrpcPort string
	dbHost         string
	dbPort         string
	dbUser         string
	dbPassword     string
	dbName         string
	redisHost      string
	redisPort      string
}

// Load config from environment variables
func loadConfig() config {
	var conf = config{
		httpPort:       env.GetEnv("HTTP_PORT"),
		authGrpcHost:   env.GetEnv("AUTH_GRPC_HOST"),
		authGrpcPort:   env.GetEnv("AUTH_GRPC_PORT"),
		thingsGrpcHost: env.GetEnv("THINGS_GRPC_HOST"),
		thingsGrpcPort: env.GetEnv("THINGS_GRPC_PORT"),
		dbHost:         env.GetEnv("PG_HOST"),
		dbPort:         env.GetEnv("PG_PORT"),
		dbName:         env.GetEnv("PG_DATABASE"),
		dbUser:         env.GetEnv("PG_USER"),
		dbPassword:     env.GetEnv("PG_PASSWORD"),
		redisHost:      env.GetEnv("REDIS_HOST"),
		redisPort:      env.GetEnv("REDIS_PORT"),
	}
	return conf
}

// createService creates a new Thing service
func createService() (service.Things, error) {
	conf := loadConfig()
	dbConf := db.DbConfig{
		Host:     conf.dbHost,
		Port:     conf.dbPort,
		Dbname:   conf.dbName,
		User:     conf.dbUser,
		Password: conf.dbPassword,
	}
	database, err := db.Init(dbConf)
	if err != nil {
		debug.LogError("Database Init error:", err)
	}
	thingsRepo := db.New(database)
	authClientConf := authClient.Conf{
		Host: conf.authGrpcHost,
		Port: conf.authGrpcPort,
	}
	authClient := authClient.New(authClientConf)
	err = authClient.StartAuthClient()
	if err != nil {
		debug.LogError("Auth client error:", err)
		return nil, err
	}
	cache := cdb.New(cdb.Conf{
		RedisHost: conf.redisHost,
		RedisPort: conf.redisPort,
	})
	err = cache.Connect()
	if err != nil {
		log.Panic(err)
	}
	return service.New(&thingsRepo, &authClient, cache), nil
}

// startHttp starts the HTTP server
func startHttp(s service.Things) error {
	conf := loadConfig()
	err := http.InitApi(conf.httpPort, s)
	return err
}
