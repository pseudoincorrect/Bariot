package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/nats-io/nats.go"
	"github.com/pseudoincorrect/bariot/pkg/env"

	"github.com/pseudoincorrect/bariot/pkg/errors"
)

func main() {
	const natsThingsSubject = "thingsMsg.>"
	const natsThingsQueue = "things"
	config := loadConfig()
	// _, err := connectToInfluxdb(config)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// log.Printf("Connected to InfluxDB")
	opts := []nats.Option{nats.Name("NATS Sample Queue Subscriber")}
	opts = natsSetupConnOptions(opts)
	natsConn, err := natsConnect(config, opts)
	if err != nil {
		log.Panic(err)
	}
	// defer natsDisconnect(natsConn)
	log.Printf("Connected to nats %s", natsConn.ConnectedUrl())

	err = natsSubscribe(natsConn, natsThingsSubject, natsThingsQueue, natsThingMsgHandler)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("Subscribed to NATS", natsThingsSubject)

	time.Sleep(20 * time.Second)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Printf("Draining...")
	natsConn.Drain()
	log.Fatalf("Exiting")

}

type config struct {
	bariotEnv      string
	influxdbHost   string
	influxdbort    string
	influxdbOrg    string
	influxdbBucket string
	influxdbToken  string
	natsHost       string
	natsPort       string
}

func loadConfig() config {
	var conf = config{
		bariotEnv:      env.GetEnv("BARIOT_ENV"),
		influxdbHost:   env.GetEnv("INFLUXDB_HOST"),
		influxdbort:    env.GetEnv("INFLUXDB_PORT"),
		influxdbOrg:    env.GetEnv("INFLUXDB_ORG"),
		influxdbBucket: env.GetEnv("INFLUXDB_BUCKET"),
		influxdbToken:  env.GetEnv("INFLUXDB_TOKEN"),
		natsHost:       env.GetEnv("NATS_HOST"),
		natsPort:       env.GetEnv("NATS_PORT"),
	}

	return conf
}

func connectToInfluxdb(cfg config) (influxdb.Client, error) {
	dbUrl := fmt.Sprintf("http://%s:%s", cfg.influxdbHost, cfg.influxdbort)
	client := influxdb.NewClient(dbUrl, cfg.influxdbToken)
	_, err := client.Health(context.Background())
	return client, err
}

func natsConnect(cfg config, opts []nats.Option) (*nats.Conn, error) {
	natsUrl := "nats://" + cfg.natsHost + ":" + cfg.natsPort
	log.Printf("Connecting to NATS Server: %s", natsUrl)
	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		return nil, errors.ErrConnection
	}
	return nc, nil
}

// func natsDisconnect(nc *nats.Conn) {
// 	nc.Close()
// }

func natsSetupConnOptions(opts []nats.Option) []nats.Option {
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

func natsSubscribe(nc *nats.Conn, subject string, queue string, handler nats.MsgHandler) error {
	nc.QueueSubscribe(subject, queue, handler)
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}

	return nil
}

func natsThingMsgHandler(msg *nats.Msg) {
	printNatsMsg(msg, 0)
}

func printNatsMsg(m *nats.Msg, i int) {
	log.Printf("Nats Message Received on [%s] Queue[%s] Pid[%d]: '%s'", m.Subject, m.Sub.Queue, os.Getpid(), string(m.Data))
}
