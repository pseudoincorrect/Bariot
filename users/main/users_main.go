package main

import (
	"context"
	"log"
	"time"

	authClient "github.com/pseudoincorrect/bariot/pkg/auth/client"
	"github.com/pseudoincorrect/bariot/pkg/env"
	"github.com/pseudoincorrect/bariot/users/api"
	"github.com/pseudoincorrect/bariot/users/db"
	"github.com/pseudoincorrect/bariot/users/models"
	"github.com/pseudoincorrect/bariot/users/service"
	"github.com/pseudoincorrect/bariot/users/utilities/hash"
)

func main() {
	log.Println("Users service online")
	usersService, err := createService()
	if err != nil {
		log.Panic("Users service creation error", err)
	}
	err = createAdmin(usersService)
	if err != nil {
		log.Panic("Admin creation error", err)
	}
	log.Println("init user service HTTP server")
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
	rpcAuthPort   string
	rpcAuthHost   string
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
		rpcAuthPort:   env.GetEnv("RPC_AUTH_PORT"),
		rpcAuthHost:   env.GetEnv("RPC_AUTH_HOST"),
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
		log.Println("Database Init error:", err)
		return nil, err
	}
	usersRepo := db.New(database)
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
	return service.New(usersRepo, authClient), nil

}

func startHttp(s service.Users) error {
	conf := loadConfig()
	return api.InitApi(conf.httpPort, s)
}

func createAdmin(s service.Users) error {
	conf := loadConfig()
	user, err := s.GetByEmail(context.Background(), conf.adminEmail)
	if err != nil {
		return err
	}
	if user == nil {
		log.Println("Admin does not exist, creating him...")
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
		_, err = s.SaveUser(context.Background(), admin)
		if err != nil {
			return err
		}
	}
	return nil
}
