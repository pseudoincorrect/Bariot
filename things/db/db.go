package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type Database struct {
	c *pgx.Conn
}

type ThingModel struct {
	Email string
}

func CreateUser(db Database, email string) (ThingModel, error) {
	var thing ThingModel
	thing.Email = email
	// err := db.Create(&thing).Error
	return thing, nil
}

type DbConfig struct {
	host     string
	user     string
	password string
	dbname   string
	port     string
}

func Init(conf DbConfig) (*Database, error) {
	db, err := connect(conf)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connect(conf DbConfig) (*Database, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.user, conf.password, conf.host, conf.port, conf.dbname)

	conn, err := pgx.Connect(context.Background(), dbUrl)

	if err != nil {
		fmt.Println("Unable to connect to database: %v\n", err)
		return nil, err
	}
	return &Database{conn}, nil
}
