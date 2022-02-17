package main

import (
	"context"
	"fmt"
	"time"

	"github.com/pseudoincorrect/bariot/users/api"
	"github.com/pseudoincorrect/bariot/users/db"
	"github.com/pseudoincorrect/bariot/users/service"
	util "github.com/pseudoincorrect/bariot/users/utilities"

	pb "github.com/pseudoincorrect/bariot/users/rpc/auth"
	"google.golang.org/grpc"
)

type config struct {
	httpPort    string
	rpcAuthPort string
	rpcAuthHost string
	dbHost      string
	dbPort      string
	dbUser      string
	dbPassword  string
	dbName      string
}

func loadConfig() config {
	var conf = config{
		httpPort:    util.GetEnv("HTTP_PORT"),
		rpcAuthPort: util.GetEnv("RPC_AUTH_PORT"),
		rpcAuthHost: util.GetEnv("RPC_AUTH_HOST"),
		dbHost:      util.GetEnv("PG_HOST"),
		dbPort:      util.GetEnv("PG_PORT"),
		dbName:      util.GetEnv("PG_DATABASE"),
		dbUser:      util.GetEnv("PG_USER"),
		dbPassword:  util.GetEnv("PG_PASSWORD"),
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
	conf := loadConfig()

	usersService := createService()
	fmt.Println("init user service HTTP server")

	go func() {
		testHttp(usersService)
	}()

	// time.Sleep(time.Second * 8)

	addr := conf.rpcAuthHost + ":" + conf.rpcAuthPort

	fmt.Println("init user service GRPC client to ", addr)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		fmt.Println("did not connect:", err)
	}
	defer conn.Close()
	c := pb.NewAuthClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	fmt.Println("GRPC get admin token request")

	res, err := c.GetAdminToken(ctx, &pb.GetAdminTokenRequest{Password: "admin_password"})

	if err != nil {
		fmt.Println("GRPC get admin token error:", err)
	}
	if err == nil {
		fmt.Println(res.GetJwt())
	}

	fmt.Println("shutting down Users service")

}
