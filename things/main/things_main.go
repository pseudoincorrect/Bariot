package main

import (
	"log"

	authClient "github.com/pseudoincorrect/bariot/pkg/auth/client"
	cdb "github.com/pseudoincorrect/bariot/pkg/cache"
	"github.com/pseudoincorrect/bariot/pkg/env"
	"github.com/pseudoincorrect/bariot/things/api"
	"github.com/pseudoincorrect/bariot/things/db"
	"github.com/pseudoincorrect/bariot/things/service"
)

func main() {
	log.Println("Things service online")
	thingsService, err := createService()
	if err != nil {
		log.Panic("Create service error:", err)
	}
	err = startHttp(thingsService)
	if err != nil {
		log.Panic("Start http error:", err)
	}
}

type config struct {
	httpPort    string
	rpcAuthHost string
	rpcAuthPort string
	dbHost      string
	dbPort      string
	dbUser      string
	dbPassword  string
	dbName      string
	redisHost   string
	redisPort   string
}

// Load config from environment variables
func loadConfig() config {
	var conf = config{
		httpPort:    env.GetEnv("HTTP_PORT"),
		rpcAuthHost: env.GetEnv("RPC_AUTH_HOST"),
		rpcAuthPort: env.GetEnv("RPC_AUTH_PORT"),
		dbHost:      env.GetEnv("PG_HOST"),
		dbPort:      env.GetEnv("PG_PORT"),
		dbName:      env.GetEnv("PG_DATABASE"),
		dbUser:      env.GetEnv("PG_USER"),
		dbPassword:  env.GetEnv("PG_PASSWORD"),
		redisHost:   env.GetEnv("REDIS_HOST"),
		redisPort:   env.GetEnv("REDIS_PORT"),
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
		log.Println("Database Init error:", err)
	}
	thingsRepo := db.New(database)

	authClientConf := authClient.Conf{
		Host: conf.rpcAuthHost,
		Port: conf.rpcAuthPort,
	}

	authClient := authClient.New(authClientConf)
	err = authClient.StartAuthClient()
	if err != nil {
		log.Println("Auth client error:", err)
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

	return service.New(thingsRepo, authClient, cache), nil
}

// startHttp starts the HTTP server
func startHttp(s service.Things) error {
	conf := loadConfig()
	err := api.InitApi(conf.httpPort, s)
	return err
}
