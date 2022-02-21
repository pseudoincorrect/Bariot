package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pseudoincorrect/bariot/users/api"
	"github.com/pseudoincorrect/bariot/users/db"
	"github.com/pseudoincorrect/bariot/users/models"
	"github.com/pseudoincorrect/bariot/users/rpc/client"
	"github.com/pseudoincorrect/bariot/users/service"
	util "github.com/pseudoincorrect/bariot/users/utilities"
)

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
		httpPort:      util.GetEnv("HTTP_PORT"),
		rpcAuthPort:   util.GetEnv("RPC_AUTH_PORT"),
		rpcAuthHost:   util.GetEnv("RPC_AUTH_HOST"),
		dbHost:        util.GetEnv("PG_HOST"),
		dbPort:        util.GetEnv("PG_PORT"),
		dbName:        util.GetEnv("PG_DATABASE"),
		dbUser:        util.GetEnv("PG_USER"),
		dbPassword:    util.GetEnv("PG_PASSWORD"),
		adminEmail:    util.GetEnv("ADMIN_EMAIL"),
		adminPassword: util.GetEnv("ADMIN_PASSWORD"),
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
		fmt.Println("Database Init error:", err)
	}
	usersRepo := db.New(database)

	authClientConf := client.AuthClientConf{
		Host: conf.rpcAuthHost,
		Port: conf.rpcAuthPort,
	}

	authClient := client.New(authClientConf)
	err = authClient.StartAuthClient()
	if err != nil {
		fmt.Println("Auth client error:", err)
		return nil, err
	}

	return service.New(usersRepo, authClient), nil

}

func testHttp(s service.Users) {
	conf := loadConfig()
	api.InitApi(conf.httpPort, s)
}

func createAdmin(s service.Users) error {
	conf := loadConfig()
	user, err := s.GetByEmail(context.Background(), conf.adminEmail)
	if err != nil {
		return err
	}
	if user == nil {
		fmt.Println("Admin does not exist, creating him...")
		hashPass, err := util.HashPassword(conf.adminPassword)
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

func main() {
	fmt.Println("Users service online")

	usersService, err := createService()
	if err != nil {
		fmt.Println("Service error:", err)
		return
	}

	createAdmin(usersService)

	fmt.Println("init user service HTTP server")

	go func() {
		testHttp(usersService)
	}()

	for {
		time.Sleep(time.Second)
	}
	// jwt, err := userService..GetUserToken(context.Background(), "123456789")
	// fmt.Println("jwt:", jwt)

	// isUser, userId, err := authClient.IsWhichUser(context.Background(), jwt)
	// if err != nil {
	// 	fmt.Println("GRPC validate token error:", err)
	// }
	// if err == nil {
	// 	fmt.Println("isUser:", isUser, ", UserId: ", userId)
	// }
}
