package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mainflux/senml"
	"github.com/nats-io/nats.go"
	"github.com/pseudoincorrect/bariot/pkg/env"
	"github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils"
	"github.com/pseudoincorrect/bariot/pkg/writer"
)

func main() {
	var w influxdbWriter
	const natsThingsSubject = "thingsMsg.>"
	const natsThingsQueue = "things"
	config := loadConfig()
	w.conf = &config
	err := w.connectToInfluxdb()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Connected to InfluxDB")
	opts := []nats.Option{nats.Name("NATS Sample Queue Subscriber")}
	opts = natsSetupConnOptions(opts)
	err = w.natsConnect(opts)
	if err != nil {
		log.Panic(err)
	}
	defer w.natsDisconnect()
	log.Println("Connected to nats", w.natsConn.ConnectedUrl())

	err = w.natsSubscribe(natsThingsSubject, natsThingsQueue, w.getNatsMsgHandler())
	if err != nil {
		log.Panic(err)
	}
	log.Println("Subscribed to NATS", natsThingsSubject)
	time.Sleep(20 * time.Second)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println()
	log.Println("Draining...")
	w.natsConn.Drain()
	log.Fatalf("Exiting")

}

type influxdbWriter struct {
	influxClient influxdb.Client
	natsConn     *nats.Conn
	conf         *config
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

// loadConfig load variable from env variables
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

// connectToInfluxdb setup a connection to an influxdb and check for health
func (w *influxdbWriter) connectToInfluxdb() error {
	dbUrl := fmt.Sprintf("http://%s:%s", w.conf.influxdbHost, w.conf.influxdbort)
	client := influxdb.NewClientWithOptions(dbUrl, w.conf.influxdbToken, influxdb.DefaultOptions().SetBatchSize(2))
	_, err := client.Health(context.Background())
	if err != nil {
		log.Panic("could not connect to influxdb")
	}
	w.influxClient = client
	return nil
}

// influxdbWrite Write a batch of senml msg to influxdb
func (w *influxdbWriter) influxdbWrite(data *writer.ThingData) {
	writeAPI := w.influxClient.WriteAPI(w.conf.influxdbOrg, w.conf.influxdbBucket)
	errChan := writeAPI.Errors()

	go func() {
		for err := range errChan {
			log.Println("Influxdb write error: ", err)
		}
	}()

	for _, r := range data.SensorsData.Records {
		p := influxdb.NewPointWithMeasurement(r.Name).
			AddTag("unit", r.Unit).
			AddField("thingId", data.ThingId).
			SetTime(time.Unix(int64(r.Time), 0))
		if r.Value != nil {
			p.AddField("value", r.Value)
		} else if r.StringValue != nil {
			p.AddField("value", r.StringValue)
		} else if r.BoolValue != nil {
			p.AddField("value", r.BoolValue)
		} else {
			log.Println("No value found in senml")
			continue
		}
		writeAPI.WritePoint(p)
	}
	// Flush writes
	writeAPI.Flush()
}

// natsConnect setup a connection to a Nats server
// TODO: check for Nats health
func (w *influxdbWriter) natsConnect(opts []nats.Option) error {
	natsUrl := "nats://" + w.conf.natsHost + ":" + w.conf.natsPort
	log.Println("Connecting to NATS Server:", natsUrl)
	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		return errors.ErrConnection
	}
	w.natsConn = nc
	return nil
}

// natsDisconnect disconnect from current nats server
func (w *influxdbWriter) natsDisconnect() {
	w.natsConn.Close()
}

// natsSetupConnOptions setup the Nats connection option such as reconnect
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
		log.Panic("Exiting:", nc.LastError())
	}))
	return opts
}

// natsSubscribe subscribe to a topic/subject with a custom handler and a queue
func (w *influxdbWriter) natsSubscribe(subject string, queue string, handler nats.MsgHandler) error {
	w.natsConn.QueueSubscribe(subject, queue, handler)
	w.natsConn.Flush()
	if err := w.natsConn.LastError(); err != nil {
		log.Panic(err)
	}
	return nil
}

// natsThingsMsgHandler handles nats message by decoding to senml format
func (w *influxdbWriter) getNatsMsgHandler() nats.MsgHandler {
	return func(natsMsg *nats.Msg) {
		printNatsMsg(natsMsg)
		msg, err := decodeNatsThingMsg(natsMsg)
		if err != nil {
			return
		}
		w.influxdbWrite(msg)
	}
}

// printNatsMsg print a nats message
func printNatsMsg(m *nats.Msg) {
	msg, err := utils.PrettyJsonString(string(m.Data))
	if err != nil {
		log.Println("Error printing Nats message")
		return
	}
	log.Printf("NATS Message Received on [%s] Queue[%s] Pid[%d]", m.Subject, m.Sub.Queue, os.Getpid())
	log.Printf("NATS Message Payload %s", msg)
}

// decodeSenmlMsg decodes a JSON message into a SenML message
func decodeNatsThingMsg(msg *nats.Msg) (*writer.ThingData, error) {
	senmlMsg, err := senml.Decode(msg.Data, senml.JSON)
	if err != nil {
		log.Println("Error decoding SenML message:", err)
		return nil, errors.ErrValidation
	}
	senmlMsg, err = senml.Normalize(senmlMsg)
	if err != nil {
		log.Println("Error normalizing SenML message:", err)
		return nil, errors.ErrValidation
	}
	for _, senmlRecord := range senmlMsg.Records {
		log.Println(utils.PrettySenmlRecord(senmlRecord))
	}
	// log.Println("Decoded SenML message:", senmlMsg)
	thingId, err := getThingIdFromNatsSubjet(msg.Subject)
	if err != nil {
		return nil, errors.ErrValidation
	}
	thingData := writer.ThingData{
		ThingId:     thingId,
		SensorsData: &senmlMsg,
	}
	return &thingData, nil
}

// getThingIdFromNatsSubjet extract the
func getThingIdFromNatsSubjet(subjet string) (string, error) {
	splits := strings.Split(subjet, ".")
	return splits[len(splits)-1], nil
}
