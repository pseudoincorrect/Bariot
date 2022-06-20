package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
)

type Database struct {
	conn *pgx.Conn
}

type DbConfig struct {
	Host     string
	Port     string
	Dbname   string
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

// connect creates a new database connection
func connect(conf DbConfig) (*Database, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Dbname)
	conn, err := pgx.Connect(context.Background(), dbUrl)

	if err != nil {
		logger.Error("Unable to connect to database:", err)
		return nil, err
	}
	return &Database{conn}, nil
}

func createThingsTable(db *Database) error {
	createTable := `create table things (
		id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		created_at timestamp with time zone DEFAULT now(),
		key character varying(4096) NOT NULL,
		CHECK (key <> ''),
		name character varying(255) NOT NULL,
		CHECK (name <> ''),
		user_id uuid NOT NULL,
		metadata json
	);`

	_, err := db.conn.Exec(context.Background(), createTable)
	if err != nil {
		logger.Error("Unable to begin :", err)
		return err
	}
	return nil
}
