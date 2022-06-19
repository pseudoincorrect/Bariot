package nats

import (
	"log"
	"time"

	natsClient "github.com/nats-io/nats.go"
	e "github.com/pseudoincorrect/bariot/pkg/errors"
)

type PubSub interface {
	Connect(opts []natsClient.Option) error
	Subscribe(subject string, queue string, handler natsClient.MsgHandler) error
}

type NatsPubSub struct {
	host string
	port string
	conn *natsClient.Conn
}

var _ PubSub = (*NatsPubSub)(nil)

func New(host string, port string) NatsPubSub {
	return NatsPubSub{
		host: host,
		port: port,
		conn: nil,
	}
}

// connect setup a connection to a Nats server
func (nps *NatsPubSub) Connect(opts []natsClient.Option) error {
	natsUrl := "nats://" + nps.host + ":" + nps.port
	log.Println("Connecting to NATS Server:", natsUrl)
	nc, err := natsClient.Connect(natsUrl, opts...)
	if err != nil {
		return e.Handle(e.ErrConn, err, "nats connect")
	}
	nps.conn = nc
	return nil
}

// natsSubscribe subscribe to a topic/subject with a custom handler and a queue
func (nps *NatsPubSub) Subscribe(subject string, queue string, handler natsClient.MsgHandler) error {
	nps.conn.QueueSubscribe(subject, queue, handler)
	nps.conn.Flush()
	if err := nps.conn.LastError(); err != nil {
		log.Panic(err)
	}
	return nil
}

// Publish publishes a message to the NATS server.
func (nps *NatsPubSub) Publish(subject string, payload string) error {
	if err := nps.conn.Publish(subject, []byte(payload)); err != nil {
		log.Printf("Error publishing message: %s\n", err)
		return err
	}
	return nil
}

// natsSetupConnOptions setup the Nats connection option such as reconnect
func NatsSetupConnOptions(name string) []natsClient.Option {
	opts := []natsClient.Option{natsClient.Name(name)}
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
		log.Panic("Exiting:", nc.LastError())
	}))
	return opts
}
