package mqtt

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/mainflux/senml"
	nats "github.com/pseudoincorrect/bariot/pkg/nats/client"
	e "github.com/pseudoincorrect/bariot/pkg/utils/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/logger"
)

type MqttSub interface {
	Connect() error
	Subscriber(topic string, qos byte, authorizer Authorizer, handler nats.NatsPubType) error
	Unsubscribe(topic string)
	Disconnect()
}

// Static type checking
var _ MqttSub = (*mqttSub)(nil)

func New(config Conf) MqttSub {
	return &mqttSub{client: nil, conf: config}
}

type mqttSub struct {
	client paho.Client
	conf   Conf
}

type Conf struct {
	User       string
	Pass       string
	Host       string
	Port       string
	HealthPort string
}

// healthCheckBlocking blocks until the health check returns a 200 OK response.
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

// healthCheck sends a health check request to the MQTT broker.
func (sub *mqttSub) healthCheck() error {
	url := "http://" + sub.conf.User + ":" + sub.conf.Pass + "@" +
		sub.conf.Host + ":" + sub.conf.HealthPort + "/api/v4/brokers"
	resp, err := http.Get(url)
	if err != nil {
		return e.ErrConn
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return e.ErrConn
	}
	return nil
}

// Connect connects to the MQTT broker.
func (sub *mqttSub) Connect() error {
	err := sub.healthCheckBlocking()
	if err != nil {
		return e.ErrConn
	}
	// paho.DEBUG = log.New(os.Stdout, "", 0)
	paho.WARN = log.New(os.Stdout, "", 0)
	paho.ERROR = log.New(os.Stdout, "", 0)
	r1 := rand.New(rand.NewSource(time.Now().UnixNano()))
	clientId := "bariot_" + strconv.Itoa(r1.Intn(1000000))
	logger.Info("MQTT client ID :", clientId)
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
	sub.client = c
	return nil
}

// Subscriber subscribes to the specified topic with a function for authorizing the topic
// and a function to handle the messages.
func (sub *mqttSub) Subscriber(topic string, qos byte,
	authorizer Authorizer, handler nats.NatsPubType) error {

	stringHandler := func(client paho.Client, msg paho.Message) {
		msgTopic := msg.Topic()
		msgPayload := msg.Payload()
		// logger.Debug("MQTT msg topic:  ", msgTopic)
		// logger.Debug("MQTT payload:", string(msgPayload))
		jwt, sensorData, err := ExtractData(msgPayload)
		if err != nil {
			logger.Error(err.Error())
		}
		err = authorizer(msgTopic, jwt)
		if err != nil {
			logger.Error(err.Error())
		}
		msgSensors, _ := senml.Encode(sensorData, senml.JSON)

		splits := strings.Split(msgTopic, "/")
		thingId := splits[len(splits)-1]
		handler(thingId, string(msgSensors))
	}

	token := sub.client.Subscribe(topic, qos, stringHandler)

	if token.Wait() && token.Error() != nil {
		log.Panic(token.Error())
	}
	return nil
}

// Unsubscribe unsubscribes from the specified topic.
func (sub *mqttSub) Unsubscribe(topic string) {
	token := sub.client.Unsubscribe(topic)
	if token.Wait() && token.Error() != nil {
		e.HandleFatal(e.ErrMqtt, nil, "Error unsubscribing from topic")
	}
}

// Disconnect disconnects from the MQTT broker.
func (sub *mqttSub) Disconnect() {
	sub.client.Disconnect(250)
}

// defaultMessageHandler handles the messages received from the MQTT broker.
var defaultMessageHandler paho.MessageHandler = func(client paho.Client, msg paho.Message) {
	logger.Debug("INCORRECT PUBLISH HERE: ", msg.Topic())
	logger.Debug("MSG: ", msg.Payload())
}

type AuthenticatedMsg struct {
	Token   string `json:"token"`
	Records []senml.Record
}

// ExtractData extracts/parse the data from the MQTT message.
func ExtractData(payload []byte) (string, senml.Pack, error) {
	msg := AuthenticatedMsg{}

	err := json.Unmarshal(payload, &msg)
	if err != nil {
		err := e.Handle(e.ErrParsing, err, "json unmarshal extract data")
		return "", senml.Pack{}, err
	}
	// logger.("JSON decoded jwt = ", msg.Token)
	// logger.("JSON decoded data = ", msg.Sensors)
	pack := senml.Pack{Records: msg.Records}
	return msg.Token, pack, nil
}
