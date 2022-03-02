package main

import (
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/nats-io/nats.go"
	"github.com/pseudoincorrect/bariot/pkg/env"
	"github.com/pseudoincorrect/bariot/pkg/errors"
)

func main() {
	config := loadConfig()

	mqttClient, err := mqttConnect(config)
	if err != nil {
		log.Panic(err)
	}
	defer mqttDisconnect(mqttClient)

	log.Printf("Connected to MQTT broker %s:%s\n", config.mqttHost, config.mqttPort)

	opts := []nats.Option{nats.Name("NATS Sample Queue Subscriber")}
	opts = natsSetupConnOptions(opts)

	natsConn, err := natsConnect(config, opts)
	if err != nil {
		log.Panic(err)
	}
	// defer natsDisconnect(natsConn)
	log.Printf("Connected to nats %s", natsConn.ConnectedUrl())

	const mqttThingsTopic = "things/#"
	const natsThingsSubject = "thingsMsg.>"

	natsPub := natsPublisher(natsConn, natsThingsSubject)

	err = mqttSubscriber(mqttClient, mqttThingsTopic, 0, natsPub)
	if err != nil {
		log.Panic(err)
	}
	defer mqttClient.Disconnect(250)
	defer mqttUnsubscribe(mqttClient, mqttThingsTopic)

	for {
		time.Sleep(15 * time.Second)
	}
}

type config struct {
	bariotEnv string
	mqttHost  string
	mqttPort  string
	natsHost  string
	natsPort  string
}

func loadConfig() config {
	var conf = config{
		bariotEnv: env.GetEnv("BARIOT_ENV"),
		mqttHost:  env.GetEnv("MQTT_HOST"),
		mqttPort:  env.GetEnv("MQTT_PORT"),
		natsHost:  env.GetEnv("NATS_HOST"),
		natsPort:  env.GetEnv("NATS_PORT"),
	}
	return conf
}

var defaultMessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("INCORRECT PUBLISH HERE: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
}

func mqttConnect(conf config) (mqtt.Client, error) {
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + conf.mqttHost + ":" + conf.mqttPort).SetClientID("bariot_mqtt_things")

	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetPingTimeout(1 * time.Minute)

	c := mqtt.NewClient(opts)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return c, nil
}

func mqttSubscriber(client mqtt.Client, topic string, qos byte, natsPub natsPubType) error {
	stringHandler := func(client mqtt.Client, msg mqtt.Message) {
		msgTopic := msg.Topic()
		msgPayload := msg.Payload()
		log.Printf("Got MQTT msg, topic: %s, payload %s\n", msgTopic, msgPayload)

		natsPub(string(msgPayload))
	}

	token := client.Subscribe(topic, qos, stringHandler)
	if token.Wait() && token.Error() != nil {
		log.Panic(token.Error())
	}
	return nil
}

func mqttUnsubscribe(client mqtt.Client, topic string) {
	token := client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Error unsubscribing from topic: %s\n", token.Error())
	}
}

func mqttDisconnect(client mqtt.Client) {
	client.Disconnect(250)
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

func natsConnect(cfg config, opts []nats.Option) (*nats.Conn, error) {
	natsUrl := "nats://" + cfg.natsHost + ":" + cfg.natsPort
	nc, err := nats.Connect(natsUrl, opts...)
	if err != nil {
		return nil, errors.ErrConnection
	}
	return nc, nil
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

func natsPublisher(nc *nats.Conn, subject string) natsPubType {
	return func(payload string) error {
		return natsPublish(nc, subject, payload)
	}
}
