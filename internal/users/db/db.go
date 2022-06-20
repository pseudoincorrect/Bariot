package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/pseudoincorrect/bariot/pkg/utils/debug"
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
		debug.LogError("Unable to connect to database:", err)
		return nil, err
	}
	return &Database{conn}, nil
}

func createUsersTable(db *Database) error {
	createTable := `create table users (
		id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
		created_at timestamp with time zone DEFAULT now(),
		email character varying(255) NOT NULL UNIQUE,
		CHECK (email <> ''),
		full_name character varying(255) NOT NULL,
		CHECK (full_name <> ''),
		hash_pass character varying(255) NOT NULL,
		CHECK (hash_pass <> ''),
		metadata json
	);`

	_, err := db.conn.Exec(context.Background(), createTable)
	if err != nil {
		debug.LogError("Unable to begin :", err)
		return err
	}
	return nil
}
