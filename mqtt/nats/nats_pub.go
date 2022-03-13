package natsPub

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pseudoincorrect/bariot/pkg/errors"
)

type NatsPub interface {
	Connect(host string, port string) error
	Disconnect()
	Publish(subject string, payload string) error
	CreatePublisher(subject string) NatsPubType
}

var _ NatsPub = (*natsPub)(nil)

func New() NatsPub {
	return &natsPub{c: nil}
}

type natsPub struct {
	c *nats.Conn
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second
	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

func (conn *natsPub) Connect(host string, port string) error {
	opts := []nats.Option{nats.Name("NATS Sample Queue Subscriber")}
	opts = setupConnOptions(opts)
	natsUrl := "nats://" + host + ":" + port
	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		return errors.ErrConnection
	}
	conn.c = nc
	return nil
}

func (conn *natsPub) Disconnect() {
	conn.c.Close()
}

func (conn *natsPub) Publish(subject string, payload string) error {
	if err := conn.c.Publish(subject, []byte(payload)); err != nil {
		log.Printf("Error publishing message: %s\n", err)
		return err
	}
	return nil
}

type NatsPubType func(payload string) error

func (conn *natsPub) CreatePublisher(subject string) NatsPubType {
	return func(payload string) error {
		err := conn.Publish(subject, payload)
		return err
	}
}
