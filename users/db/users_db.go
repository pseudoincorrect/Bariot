package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v4"
)

type Database struct {
	conn *pgx.Conn
}

type DbConfig struct {
	Host     string
	Dbname   string
	Port     string
	User     string
	Password string
}

// Init a new database connection
func Init(conf DbConfig) (*Database, error) {
	db, err := connect(conf)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// connect create a new database connection
func connect(conf DbConfig) (*Database, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Dbname)
	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Println("Unable to connect to database:", err)
		return nil, err
	}
	return &Database{conn}, nil
}
