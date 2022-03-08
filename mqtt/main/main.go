package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nats-io/nats.go"
	"github.com/pseudoincorrect/bariot/pkg/env"
	"github.com/pseudoincorrect/bariot/pkg/errors"
)

func main() {
	log.SetOutput(os.Stdout)

	var f forwarder
	f.conf = loadConfig()

	err := f.mqttConnect()
	if err != nil {
		log.Panic(err)
	}
	defer f.mqttDisconnect()

	log.Printf("Connected to MQTT broker %s:%s\n", f.conf.mqttHost, f.conf.mqttPort)

	opts := []nats.Option{nats.Name("NATS Sample Queue Subscriber")}
	opts = natsSetupConnOptions(opts)

	err = f.natsConnect(opts)
	if err != nil {
		log.Panic(err)
	}
	// defer natsDisconnect(natsConn)
	log.Printf("Connected to nats %s", f.natsConn.ConnectedUrl())

	const mqttThingsTopic = "things/#"
	const natsThingsSubject = "thingsMsg.>"

	natsPub := f.createNatsPublisher(natsThingsSubject)

	err = f.mqttSubscriber(mqttThingsTopic, 0, natsPub)

	if err != nil {
		log.Panic(err)
	}
	defer f.mqttClient.Disconnect(250)
	defer f.mqttUnsubscribe(mqttThingsTopic)

	for {
		time.Sleep(5 * time.Second)
	}
}

type forwarder struct {
	natsConn   *nats.Conn
	mqttClient mqtt.Client
	conf       config
}

type config struct {
	bariotEnv      string
	mqttHost       string
	mqttPort       string
	mqttUser       string
	mqttPass       string
	mqttHealthPort string
	natsHost       string
	natsPort       string
}

func loadConfig() config {
	var conf = config{
		bariotEnv:      env.GetEnv("BARIOT_ENV"),
		mqttHost:       env.GetEnv("MQTT_HOST"),
		mqttPort:       env.GetEnv("MQTT_PORT"),
		mqttHealthPort: env.GetEnv("MQTT_HEALTH_PORT"),
		mqttUser:       env.GetEnv("MQTT_USER"),
		mqttPass:       env.GetEnv("MQTT_PASS"),
		natsHost:       env.GetEnv("NATS_HOST"),
		natsPort:       env.GetEnv("NATS_PORT"),
	}
	return conf
}

var defaultMessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("INCORRECT PUBLISH HERE: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
}

func (f *forwarder) mqttHealthCheckBlocking() error {
	for {
		err := f.mqttHealthCheck()
		if err == nil {
			return nil
		}
		fmt.Println("MQTT broker not online, retrying later...")
		time.Sleep(5 * time.Second)
	}
}

func (f *forwarder) mqttHealthCheck() error {
	url := "http://" + f.conf.mqttUser + ":" + f.conf.mqttPass + "@" +
		f.conf.mqttHost + ":" + f.conf.mqttHealthPort + "/api/v4/brokers"
	resp, err := http.Get(url)
	if err != nil {
		return errors.ErrConnection
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.ErrConnection
	}
	return nil
}

func (f *forwarder) mqttConnect() error {
	err := f.mqttHealthCheckBlocking()
	if err != nil {
		return errors.ErrConnection
	}

	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.WARN = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	clientId := "bariot_" + strconv.Itoa(r1.Intn(1000000))
	log.Println("MQTT client ID :", clientId)
	url := "tcp://" + f.conf.mqttHost + ":" + f.conf.mqttPort

	opts := mqtt.NewClientOptions().AddBroker(url).SetClientID(clientId)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	f.mqttClient = c
	return nil
}

func (f *forwarder) mqttSubscriber(topic string, qos byte, natsPub natsPubType) error {
	stringHandler := func(client mqtt.Client, msg mqtt.Message) {
		msgTopic := msg.Topic()
		msgPayload := msg.Payload()
		log.Printf("Got MQTT msg, topic: %s, payload %s\n", msgTopic, msgPayload)
		natsPub(string(msgPayload))
	}

	token := f.mqttClient.Subscribe(topic, qos, stringHandler)
	if token.Wait() && token.Error() != nil {
		log.Panic(token.Error())
	}
	return nil
}

func (f *forwarder) mqttUnsubscribe(topic string) {
	token := f.mqttClient.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Error unsubscribing from topic: %s\n", token.Error())
	}
}

func (f *forwarder) mqttDisconnect() {
	f.mqttClient.Disconnect(250)
}

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

func (f *forwarder) natsConnect(opts []nats.Option) error {
	natsUrl := "nats://" + f.conf.natsHost + ":" + f.conf.natsPort
	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		return errors.ErrConnection
	}
	f.natsConn = nc
	return nil
}

// func natsDisconnect(nc *nats.Conn) {
// 	nc.Close()
// }

func natsPublish(nc *nats.Conn, subject string, payload string) error {
	if err := nc.Publish(subject, []byte(payload)); err != nil {
		log.Printf("Error publishing message: %s\n", err)
		return err
	}
	return nil
}

type natsPubType func(payload string) error

func (f *forwarder) createNatsPublisher(subject string) natsPubType {
	return func(payload string) error {
		err := natsPublish(f.natsConn, subject, payload)
		return err
	}
}
