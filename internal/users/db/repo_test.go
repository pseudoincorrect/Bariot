package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pseudoincorrect/bariot/internal/users/models"
	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
	"github.com/stretchr/testify/assert"
)

var db *Database

const userDbHost string = "localhost"
const userDbPort string = "5432"
const userDbName string = "user_db_name"
const userDbUser string = "user_db_user"
const userDbPassword string = "user_db_password"

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14-alpine",
		Env: []string{
			"POSTGRES_DB=" + userDbName,
			"POSTGRES_USER=" + userDbUser,
			"POSTGRES_PASSWORD=" + userDbPassword,
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		logger.Error("Open Docker Settings and select 'Expose daemon on tcp://localhost:2375 without TLS'")
		e.HandleFatal(e.ErrConn, err, "Could not start resource")
	}
	ContainerPort := resource.GetPort(userDbPort + "/tcp")
	// Exponential back-off-retry, because the application in the container
	// might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		conf := DbConfig{
			Host:     userDbHost,
			Port:     ContainerPort,
			Dbname:   userDbName,
			User:     userDbUser,
			Password: userDbPassword,
		}
		dbConnect, err := connect(conf)
		if err != nil {
			return err
		}
		db = dbConnect
		return db.conn.Ping(context.Background())
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	err = createUsersTable(db)
	if err != nil {
		log.Fatalf("Could not create user table")
	}
	code := m.Run()
	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

func TestSave(t *testing.T) {
	clearUserTable()
	repo := New(db)
	testData := []struct {
		message  string
		user     models.User
		expError error
	}{
		{
			message:  "OK, should save a user",
			user:     createMockUser("", "Harry"),
			expError: nil,
		},
		{
			message:  "should throw, empty name string",
			user:     createMockUser("", ""),
			expError: e.ErrDb,
		},
	}
	for _, d := range testData {
		ctx := context.Background()
		err := repo.Save(ctx, &d.user)
		assert.Equal(t, d.expError, err, d.message)
	}
}

func TestGetByEmail(t *testing.T) {
	clearUserTable()
	repo := New(db)
	testData := []struct {
		message    string
		user       models.User
		savedUser  models.User
		gottenUser bool
	}{
		{
			message:    "OK, should save a user",
			savedUser:  createMockUser("", "Harry"),
			user:       createMockUser("", "Harry"),
			gottenUser: true,
		},
		{
			message:    "should throw, empty name string",
			savedUser:  createMockUser("", "Harry"),
			user:       createMockUser("", "Johnny"),
			gottenUser: false,
		},
	}
	for _, d := range testData {
		clearUserTable()
		ctx := context.Background()
		err := repo.Save(ctx, &d.savedUser)
		assert.Nil(t, err, "should save user without error")
		usr, err := repo.GetByEmail(ctx, d.user.Email)
		assert.Nil(t, err, "should get user without error")
		if d.gottenUser {
			assert.Equal(t, d.user.FullName, usr.FullName, d.message)
		} else {
			assert.Nil(t, usr, d.message)
		}
	}
}

func TestGet(t *testing.T) {
	clearUserTable()
	repo := New(db)
	assert.NotNil(t, repo, "repo should be created")

}

func TestDelete(t *testing.T) {
	clearUserTable()
	repo := New(db)
	assert.NotNil(t, repo, "repo should be created")

}

func TestUpdate(t *testing.T) {
	clearUserTable()
	repo := New(db)
	assert.NotNil(t, repo, "repo should be created")

}

func clearUserTable() {
	createTable := `TRUNCATE users`
	_, err := db.conn.Exec(context.Background(), createTable)
	if err != nil {
		log.Fatal("Unable to empty user table :", err)
	}
}

func createMockUser(id string, name string) models.User {
	return models.User{
		Id:        id,
		CreatedAt: "",
		Email:     name + "@test.com",
		FullName:  name,
		HashPass:  "ldzkjhgbksdqlkmfdsljkfsmlk",
		Metadata:  models.Metadata{"address": "Cornimont"},
	}
}
