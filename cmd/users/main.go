package main

import (
	"context"
	"log"
	"time"

	"github.com/pseudoincorrect/bariot/internal/users/db"
	"github.com/pseudoincorrect/bariot/internal/users/hash"
	"github.com/pseudoincorrect/bariot/internal/users/http"
	"github.com/pseudoincorrect/bariot/internal/users/models"
	"github.com/pseudoincorrect/bariot/internal/users/service"
	authClient "github.com/pseudoincorrect/bariot/pkg/auth/client"
	"github.com/pseudoincorrect/bariot/pkg/env"
	"github.com/pseudoincorrect/bariot/pkg/utils/debug"
)

func main() {
	debug.LogInfo("Users service online")
	usersService, err := createService()
	if err != nil {
		log.Panic("Users service creation error", err)
	}
	err = createAdmin(usersService)
	if err != nil {
		log.Panic("Admin creation error", err)
	}
	debug.LogInfo("init user service HTTP server")
	go func() {
		err = startHttp(usersService)
		if err != nil {
			log.Panic("Users service HTTP server error", err)
		}
	}()
	for {
		time.Sleep(time.Second)
	}
}

type config struct {
	httpPort      string
	authGrpcPort  string
	authGrpcHost  string
	dbHost        string
	dbPort        string
	dbUser        string
	dbPassword    string
	dbName        string
	adminEmail    string
	adminPassword string
}

func loadConfig() config {
	var conf = config{
		httpPort:      env.GetEnv("HTTP_PORT"),
		authGrpcPort:  env.GetEnv("AUTH_GRPC_PORT"),
		authGrpcHost:  env.GetEnv("AUTH_GRPC_HOST"),
		dbHost:        env.GetEnv("PG_HOST"),
		dbPort:        env.GetEnv("PG_PORT"),
		dbName:        env.GetEnv("PG_DATABASE"),
		dbUser:        env.GetEnv("PG_USER"),
		dbPassword:    env.GetEnv("PG_PASSWORD"),
		adminEmail:    env.GetEnv("ADMIN_EMAIL"),
		adminPassword: env.GetEnv("ADMIN_PASSWORD"),
	}
	return conf
}

// createService creates a new User service.
func createService() (service.Users, error) {
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
		return nil, err
	}
	usersRepo := db.New(database)
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
	return service.New(&usersRepo, &authClient), nil
}

// startHttp starts the HTTP server.
func startHttp(s service.Users) error {
	conf := loadConfig()
	return http.InitApi(conf.httpPort, s)
}

// createAdmin creates a new admin user.
func createAdmin(s service.Users) error {
	conf := loadConfig()
	user, err := s.GetByEmail(context.Background(), conf.adminEmail)
	if err != nil {
		return err
	}
	if user == nil {
		debug.LogError("Admin does not exist, creating him...")
		hashPass, err := hash.HashPassword(conf.adminPassword)
		if err != nil {
			return err
		}
		admin := &models.User{
			Id:       "0",
			Email:    conf.adminEmail,
			FullName: "Admin",
			HashPass: hashPass,
		}
		err = s.SaveUser(context.Background(), admin)
		if err != nil {
			return err
		}
	}
	return nil
}
