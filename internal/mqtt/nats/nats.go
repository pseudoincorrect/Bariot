package nats

import (
	"log"
	"time"

	natsClient "github.com/nats-io/nats.go"
	"github.com/pseudoincorrect/bariot/pkg/errors"
)

type NatsPub interface {
	Connect() error
	Disconnect()
	Publish(subject string, payload string) error
	CreatePublisher(subject string) NatsPubType
}

// Static type checking
var _ NatsPub = (*natsPub)(nil)

func New(config Conf) NatsPub {
	return &natsPub{client: nil, conf: config}
}

type natsPub struct {
	client *natsClient.Conn
	conf   Conf
}

type Conf struct {
	Host string
	Port string
}

// setupConnOptions parses NATS connection options from the command line
func setupConnOptions(opts []natsClient.Option) []natsClient.Option {
	totalWait := 10 * time.Minute
	reconnectDelay := time.Second
	opts = append(opts, natsClient.ReconnectWait(reconnectDelay))
	opts = append(opts, natsClient.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, natsClient.DisconnectErrHandler(func(nc *natsClient.Conn, err error) {
		log.Printf("Disconnected due to: %s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, natsClient.ReconnectHandler(func(nc *natsClient.Conn) {
		log.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, natsClient.ClosedHandler(func(nc *natsClient.Conn) {
		log.Fatalf("Exiting: %v", nc.LastError())
	}))
	return opts
}

// Connect connects to the NATS server.
func (pub *natsPub) Connect() error {
	opts := []natsClient.Option{natsClient.Name("NATS Sample Queue Subscriber")}
	opts = setupConnOptions(opts)
	natsUrl := "nats://" + pub.conf.Host + ":" + pub.conf.Port
	nc, err := natsClient.Connect(natsUrl, opts...)
	if err != nil {
		return errors.ErrConn
	}
	pub.client = nc
	return nil
}

//Disconnect disconnects from the NATS server.
func (pub *natsPub) Disconnect() {
	pub.client.Close()
}

// Publish publishes a message to the NATS server.
func (pub *natsPub) Publish(subject string, payload string) error {
	if err := pub.client.Publish(subject, []byte(payload)); err != nil {
		log.Printf("Error publishing message: %s\n", err)
		return err
	}
	return nil
}

type NatsPubType func(thingId string, payload string) error

// CreatePublisher return a function to publish on a given subject.
func (pub *natsPub) CreatePublisher(subject string) NatsPubType {
	return func(thingId string, payload string) error {
		thingIdSubject := subject + "." + thingId
		err := pub.Publish(thingIdSubject, payload)
		return err
	}
}
