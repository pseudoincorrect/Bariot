package natsPub

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pseudoincorrect/bariot/pkg/errors"
)

type NatsPub interface {
	Connect() error
	Disconnect()
	Publish(subject string, payload string) error
	CreatePublisher(subject string) NatsPubType
}

var _ NatsPub = (*natsPub)(nil)

func New(config NatsConf) NatsPub {
	return &natsPub{c: nil, conf: config}
}

type natsPub struct {
	c    *nats.Conn
	conf NatsConf
}

type NatsConf struct {
	Host string
	Port string
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

func (pub *natsPub) Connect() error {
	opts := []nats.Option{nats.Name("NATS Sample Queue Subscriber")}
	opts = setupConnOptions(opts)
	natsUrl := "nats://" + pub.conf.Host + ":" + pub.conf.Port
	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		return errors.ErrConnection
	}
	pub.c = nc
	return nil
}

func (pub *natsPub) Disconnect() {
	pub.c.Close()
}

func (pub *natsPub) Publish(subject string, payload string) error {
	if err := pub.c.Publish(subject, []byte(payload)); err != nil {
		log.Printf("Error publishing message: %s\n", err)
		return err
	}
	return nil
}

type NatsPubType func(thingId string, payload string) error

func (pub *natsPub) CreatePublisher(subject string) NatsPubType {
	return func(thingId string, payload string) error {
		thingIdSubjet := subject + "." + thingId
		err := pub.Publish(thingIdSubjet, payload)
		return err
	}
}
