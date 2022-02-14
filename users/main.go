package main

import (
	"fmt"

	"github.com/pseudoincorrect/bariot/users/api"
	"github.com/pseudoincorrect/bariot/users/db"
	"github.com/pseudoincorrect/bariot/users/service"
	util "github.com/pseudoincorrect/bariot/users/utilities"
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
		httpPort:   util.GetEnv("HTTP_PORT"),
		dbHost:     util.GetEnv("PG_HOST"),
		dbPort:     util.GetEnv("PG_PORT"),
		dbName:     util.GetEnv("PG_DATABASE"),
		dbUser:     util.GetEnv("PG_USER"),
		dbPassword: util.GetEnv("PG_PASSWORD"),
	}
	return conf
}

func createService() service.Users {
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

	return service.New(usersRepo)

}
func testHttp(s service.Users) {
	conf := loadConfig()

	api.InitApi(conf.httpPort, s)
}

func main() {
	fmt.Println("Users service online")

	usersService := createService()

	testHttp(usersService)

}
