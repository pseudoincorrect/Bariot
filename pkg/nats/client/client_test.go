package client

import (
	"fmt"
	"os"
	"testing"

	natsGo "github.com/nats-io/nats.go"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	e "github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/debug"
)

const host = "localhost"
const port = "4222"
const subSubject = "thingsMsg.>"
const queue = "things"

var nps nats

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		e.HandleFatal(e.ErrDocker, err, "Could not connect to docker")
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
		debug.LogError("Open Docker Settings and select 'Expose daemon on tcp://localhost:2375 without TLS'")
		e.HandleFatal(e.ErrConn, err, "Could not start resource")
	}
	ContainerPort := resource.GetPort(port + "/tcp")
	nps = New(Conf{Host: host, Port: ContainerPort})
	// Exponential back-off-retry, container might not be ready
	if err := pool.Retry(func() error {
		return nps.Connect(NatsSetupConnOptions("nats reader"))
	}); err != nil {
		e.HandleFatal(e.ErrConn, err, "Could not connect to nats")
	}
	code := m.Run()
	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		e.HandleFatal(e.ErrConn, err, "Could not purge resource")
	}
	os.Exit(code)
}

func TestPubSub(t *testing.T) {
	nps.Subscribe(subSubject, queue, natsHandler)
	thingId := "000-000-001"
	pubSubject := subSubject + "." + thingId
	nps.Publish(pubSubject, "123456789")
}

func natsHandler(m *natsGo.Msg) {
	debug.LogDebug("NATS HANDLER")
	str := fmt.Sprintf("NATS Message Received on [%s] Queue[%s] Pid[%d]", m.Subject, m.Sub.Queue, os.Getpid())
	debug.LogInfo(str)
	debug.LogDebug("NATS Message Payload %s", m.Data)
}
