package main

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/pseudoincorrect/bariot/pkg/env"
)

type config struct {
	host   string
	port   string
	org    string
	bucket string
	token  string
}

func loadConfig() config {
	var conf = config{
		host:   env.GetEnv("INFLUXDB_HOST"),
		port:   env.GetEnv("INFLUXDB_PORT"),
		org:    env.GetEnv("INFLUXDB_ORG"),
		bucket: env.GetEnv("INFLUXDB_BUCKET"),
		token:  env.GetEnv("INFLUXDB_TOKEN"),
	}
	return conf
}

func connectToInfluxDB() (influxdb2.Client, error) {
	config := loadConfig()

	dbUrl := fmt.Sprintf("http://%s:%s", config.host, config.port)
	client := influxdb2.NewClient(dbUrl, config.token)

	// validate client connection health
	_, err := client.Health(context.Background())

	return client, err
}
