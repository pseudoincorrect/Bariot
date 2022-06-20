package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pseudoincorrect/bariot/internal/things/models"
	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
	"github.com/stretchr/testify/assert"
)

var db *Database

const thingDbHost string = "localhost"
const thingDbPort string = "5432"
const thingDbName string = "thing_db_name"
const thingDbUser string = "thing_db_user"
const thingDbPassword string = "thing_db_password"

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14-alpine",
		Env: []string{
			"POSTGRES_DB=" + thingDbName,
			"POSTGRES_USER=" + thingDbUser,
			"POSTGRES_PASSWORD=" + thingDbPassword,
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
	ContainerPort := resource.GetPort(thingDbPort + "/tcp")
	// Exponential back-off-retry, because the application in the container
	// might not be ready to accept connections yet
	if err := pool.Retry(func() error {
		conf := DbConfig{
			Host:     thingDbHost,
			Port:     ContainerPort,
			Dbname:   thingDbName,
			User:     thingDbUser,
			Password: thingDbPassword,
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
	err = createThingsTable(db)
	if err != nil {
		log.Fatalf("Could not create thing table")
	}
	code := m.Run()
	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestSave(t *testing.T) {
	clearThingTable()
	id := "00000000-0000-0000-0000-000000000001"
	userId := "00000000-0000-0000-0000-000000000002"
	key := "124567896"
	name := "smart_bottle_1"
	repo := New(db)
	testData := []struct {
		message  string
		thing    models.Thing
		expError error
	}{
		{
			message:  "OK, should save a thing",
			thing:    createMockThing(id, userId, key, name),
			expError: nil,
		},
		{
			message:  "should throw, wrong user ID format",
			thing:    createMockThing(id, "malformed", key, name),
			expError: e.ErrDbUuid,
		},
		{
			message:  "should throw, empty key",
			thing:    createMockThing(id, userId, "", name),
			expError: e.ErrDb,
		},
	}
	for _, d := range testData {
		ctx := context.Background()
		err := repo.Save(ctx, &d.thing)
		assert.Equal(t, d.expError, err, d.message, " | error:", err)
	}
}

func TestGet(t *testing.T) {
	clearThingTable()
	repo := New(db)
	userId := "00000000-0000-0000-0000-000000000002"
	thing := createMockThing("", userId, "key", "name")
	ctx := context.Background()
	err := repo.Save(ctx, &thing)
	assert.Nil(t, err, "saving test thing")
	testData := []struct {
		message  string
		id       string
		expError error
	}{
		{
			message:  "OK, should get the right thing",
			id:       thing.Id,
			expError: nil,
		},
		{
			message:  "Should return thing not found error",
			id:       "00000000-0000-0000-0000-000000000003",
			expError: e.ErrDbNotFound,
		},
	}
	for _, d := range testData {
		_, err := repo.Get(ctx, d.id)
		assert.Equal(t, d.expError, err, d.message, " | error:", err)
	}
}

func TestDelete(t *testing.T) {
	clearThingTable()
	repo := New(db)
	userId := "00000000-0000-0000-0000-000000000002"
	thing := createMockThing("", userId, "key", "name")
	ctx := context.Background()
	err := repo.Save(ctx, &thing)
	assert.Nil(t, err, "saving test thing")
	testData := []struct {
		message  string
		id       string
		expError error
	}{
		{
			message:  "OK, should Delete the right thing",
			id:       thing.Id,
			expError: nil,
		},
		{
			message:  "Should return thing not found error",
			id:       "00000000-0000-0000-0000-000000000003",
			expError: e.ErrDbNotFound,
		},
	}
	for _, d := range testData {
		_, err := repo.Delete(ctx, d.id)
		assert.Equal(t, d.expError, err, d.message, " | error:", err)
	}
}

func TestUpdate(t *testing.T) {
	clearThingTable()
	repo := New(db)
	userId := "00000000-0000-0000-0000-000000000002"
	name := "name"
	newName := "new_name"
	thing := createMockThing("", userId, "key", name)
	thing.Name = newName

	ctx := context.Background()
	err := repo.Save(ctx, &thing)
	assert.Nil(t, err, "saving test thing")

	err = repo.Update(ctx, &thing)
	assert.Nil(t, err, "should update without error")
	assert.Equal(t, newName, thing.Name, "should have the name updated")

	thingId := "00000000-0000-0000-0000-000000000002"
	otherThing := createMockThing(thingId, userId, "key", "name")

	err = repo.Update(ctx, &otherThing)
	assert.Equal(t, e.ErrDbNotFound, err, "should get a not found error", err)
}

func createMockThing(thingId string, userId string, key string, name string) models.Thing {
	return models.Thing{
		Id:        thingId,
		CreatedAt: "2022-06-01T14:35:40+03:00",
		Key:       key,
		Name:      name,
		UserId:    userId,
		Metadata:  models.Metadata{"unit": "temperature"},
	}
}

func clearThingTable() {
	createTable := `TRUNCATE things`
	_, err := db.conn.Exec(context.Background(), createTable)
	if err != nil {
		log.Fatal("Unable to empty thing table :", err)
	}
}
