package nats

import (
	"log"
	"os"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const host = "localhost"
const port = "4222"
const subSubject = "thingsMsg.>"
const queue = "things"

var nps NatsPubSub

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "nats",
		Tag:        "2.7-alpine",
		Env: []string{
			"NATS_HOST=" + host,
			"NATS_PORT=" + port,
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Println("Open Docker Settings and select 'Expose daemon on tcp://localhost:2375 without TLS'")
		log.Fatalf("Could not start resource: %s", err)
	}
	ContainerPort := resource.GetPort(port + "/tcp")
	nps = New(host, ContainerPort)
	// Exponential back-off-retry, container might not be ready
	if err := pool.Retry(func() error {
		return nps.Connect(NatsSetupConnOptions("nats reader"))
	}); err != nil {
		log.Fatalf("Could not connect to nats: %s", err)
	}
	code := m.Run()
	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
	os.Exit(code)
}

func TestPubSub(t *testing.T) {
	nps.Subscribe(subSubject, queue, natsHandler)
	thingId := "000-000-001"
	pubSubject := subSubject + "." + thingId
	nps.Publish(pubSubject, "123456789")
}

func natsHandler(m *nats.Msg) {
	log.Printf("NATS HANDLER")
	log.Printf("NATS Message Received on [%s] Queue[%s] Pid[%d]", m.Subject, m.Sub.Queue, os.Getpid())
	log.Printf("NATS Message Payload %s", m.Data)
}
