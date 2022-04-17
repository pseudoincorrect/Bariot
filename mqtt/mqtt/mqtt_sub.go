package mqtt

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	natsPub "github.com/pseudoincorrect/bariot/mqtt/nats"
	"github.com/pseudoincorrect/bariot/pkg/errors"
)

type MqttSub interface {
	Connect() error
	Subscriber(topic string, qos byte, authorizer Authorizer, handler natsPub.NatsPubType) error
	Unsubscribe(topic string)
	Disconnect()
}

var _ MqttSub = (*mqttSub)(nil)

func New(config MqttSubConf) MqttSub {
	return &mqttSub{c: nil, conf: config}
}

type mqttSub struct {
	c    paho.Client
	conf MqttSubConf
}

type MqttSubConf struct {
	User       string
	Pass       string
	Host       string
	Port       string
	HealthPort string
}

func (sub *mqttSub) healthCheckBlocking() error {
	for {
		err := sub.healthCheck()
		if err == nil {
			return nil
		}
		fmt.Println("MQTT broker not online, retrying later...")
		time.Sleep(5 * time.Second)
	}
}

func (sub *mqttSub) healthCheck() error {
	url := "http://" + sub.conf.User + ":" + sub.conf.Pass + "@" +
		sub.conf.Host + ":" + sub.conf.HealthPort + "/api/v4/brokers"
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

func (sub *mqttSub) Connect() error {
	err := sub.healthCheckBlocking()
	if err != nil {
		return errors.ErrConnection
	}
	// paho.DEBUG = log.New(os.Stdout, "", 0)
	paho.WARN = log.New(os.Stdout, "", 0)
	paho.ERROR = log.New(os.Stdout, "", 0)
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	clientId := "bariot_" + strconv.Itoa(r1.Intn(1000000))
	log.Println("MQTT client ID :", clientId)
	url := "tcp://" + sub.conf.Host + ":" + sub.conf.Port

	opts := paho.NewClientOptions().AddBroker(url).SetClientID(clientId)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetPingTimeout(1 * time.Second)

	c := paho.NewClient(opts)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	sub.c = c
	return nil
}

func (sub *mqttSub) Subscriber(topic string, qos byte,
	authorizer Authorizer, handler natsPub.NatsPubType) error {
	stringHandler := func(client paho.Client, msg paho.Message) {
		msgTopic := msg.Topic()
		msgPayload := msg.Payload()
		log.Printf("MQTT msg RECEIVED\n")
		log.Printf("MQTT topic:   %s\n", msgTopic)
		log.Printf("MQTT payload: %s\n", msgPayload)
		jwt, err := ExtractJwt(msgPayload)
		if err != nil {
			log.Println(err.Error())
		}
		err = authorizer(msgTopic, jwt)
		if err != nil {
			log.Println(err.Error())
		}
		handler(string(msgPayload))
	}
	token := sub.c.Subscribe(topic, qos, stringHandler)
	if token.Wait() && token.Error() != nil {
		log.Panic(token.Error())
	}
	return nil
}

func (sub *mqttSub) Unsubscribe(topic string) {
	token := sub.c.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Error unsubscribing from topic: %s\n", token.Error())
	}
}

func (sub *mqttSub) Disconnect() {
	sub.c.Disconnect(250)
}

var defaultMessageHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	log.Printf("INCORRECT PUBLISH HERE: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
}
