package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mainflux/senml"
	"github.com/pseudoincorrect/bariot/pkg/errors"
)

/// GetEnv returns the value of the environment variable named by the key.
func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Println("Environment variable", key, "is not set")
		log.Println("Please set it and try again")
		panic("Environment variable " + key + " is not set")
	}
	return value
}

type config struct {
	mqttHost       string
	mqttPort       string
	mqttUser       string
	mqttPass       string
	mqttHealthPort string
	thingToken     string
	thingId        string
}

type mqttTester struct {
	conf   config
	client mqtt.Client
}

func main() {
	MqttConnectAndSend()
}

// const BARIOT_HOST = "ec2-46-51-148-15.eu-west-1.compute.amazonaws.com"
// const THING_TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySWQiOiI5ODdiYTg5MS03NTg4LTRlMTgtOGUxOS1iMzY1NWQ1ZGJlNzYiLCJleHAiOjE2NTM1NTc3OTEsImlhdCI6MTY1MzQ3MTM5MSwiaXNzIjoiZGV2X2xvY2FsIiwic3ViIjoiMmY5NDA3MmItMDBkMy00YzZiLTk4MGYtMmY1YmY4MzJhMWNjIn0.tc0UjF3AiS965-cNAfCyccWut1255bW1w_LR3n4ZznU"
// const THING_ID = "2f94072b-00d3-4c6b-980f-2f5bf832a1cc"

// const TOPIC = "things/" + THING_ID

func MqttConnectAndSend() error {

	var m mqttTester

	mqttHost := GetEnv("BARIOT_HOST")
	thingToken := GetEnv("THING_TOKEN")
	thingId := GetEnv("THING_ID")

	var topic = "things/" + thingId

	m.conf = config{
		mqttHost:       mqttHost,
		mqttPort:       "1883",
		mqttUser:       "admin",
		mqttPass:       "public",
		mqttHealthPort: "8084",
		thingToken:     thingToken,
		thingId:        thingId,
	}
	err := m.mqttConnect()
	if err != nil {
		log.Panic("could not connect to MQTT broker")
		// return err
	}
	defer m.mqttDisconnect()
	log.Println("Connected to mqtt")
	sensorData := createSenmlPack()
	msg, _ := marchalMsg(m.conf.thingToken, sensorData)
	log.Println("Publishing to mqtt")
	err = m.mqttPublish(topic, string(msg))
	if err != nil {
		log.Panic("could not publish MQTT message")
	}
	time.Sleep(1 * time.Second)
	return nil
}

func (m *mqttTester) mqttSetOpts() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + m.conf.mqttHost + ":" + m.conf.mqttPort).SetClientID("bariot_mqtt_things")
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetPingTimeout(1 * time.Minute)
	return opts
}

func (m *mqttTester) mqttHealthCheckBlocking() error {
	for {
		err := m.mqttHealthCheck()
		if err == nil {
			return nil
		}
		fmt.Println("MQTT broker not online, retrying later...")
		time.Sleep(5 * time.Second)
	}
}

func (m *mqttTester) mqttHealthCheckOnce() error {
	err := m.mqttHealthCheck()
	if err == nil {
		return nil
	}
	return errors.ErrConnection
}

func (m *mqttTester) mqttHealthCheck() error {
	url := "http://" + m.conf.mqttUser + ":" + m.conf.mqttPass + "@" +
		m.conf.mqttHost + ":" + m.conf.mqttHealthPort + "/api/v4/brokers"
	resp, err := http.Get(url)
	if err != nil {
		print("if EMQX is behind a reverse proxy, we cannot make healthcheck")
		return errors.ErrConnection
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.ErrConnection
	}
	return nil
}

func (m *mqttTester) mqttConnect() error {
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := m.mqttSetOpts()
	// err := m.mqttHealthCheckBlocking()
	// err := m.mqttHealthCheckOnce()
	// if err != nil {
	// 	return errors.ErrConnection
	// }
	c := mqtt.NewClient(opts)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	m.client = c
	return nil
}

func (m *mqttTester) mqttDisconnect() {
	m.client.Disconnect(250)
}

func defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("INCORRECT PUBLISH HERE: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
}

func (m *mqttTester) mqttPublish(topic string, msg string) error {
	token := m.client.Publish(topic, 0, false, msg)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

type sensorAuthMsg struct {
	Token   string `json:"token"`
	Sensors senml.Pack
}

func marchalMsg(token string, sensorData senml.Pack) ([]byte, error) {
	msg := sensorAuthMsg{
		Token:   token,
		Sensors: sensorData,
	}
	return json.Marshal(msg)
}

func createSenmlPack() senml.Pack {
	temperature := float64(38)
	humidity := float64(75)
	activity := float64(1.2)
	return senml.Pack{
		Records: []senml.Record{
			{
				Name:  "temperature",
				Unit:  "degreesC",
				Value: &temperature,
				Time:  float64(time.Now().Unix()),
			},
			{
				Name:  "humidity",
				Unit:  "percents",
				Value: &humidity,
				Time:  float64(time.Now().Unix() + 1),
			},
			{
				Name:  "activity",
				Unit:  "G",
				Value: &activity,
				Time:  float64(time.Now().Unix() + 2),
			},
		},
	}
}

func createSenmlMsg() ([]byte, error) {
	temperature := float64(38)
	humidity := float64(75)
	activity := float64(1.2)
	msg := senml.Pack{
		Records: []senml.Record{
			{
				Name:  "temperature",
				Unit:  "degreesC",
				Value: &temperature,
				Time:  float64(time.Now().Unix()),
			},
			{
				Name:  "humidity",
				Unit:  "percents",
				Value: &humidity,
				Time:  float64(time.Now().Unix() + 1),
			},
			{
				Name:  "activity",
				Unit:  "G",
				Value: &activity,
				Time:  float64(time.Now().Unix() + 2),
			},
		},
	}
	enc, err := senml.Encode(msg, senml.JSON)
	if err != nil {
		log.Panic("wrong senml format during encoding")
		return nil, err
	}
	return enc, nil
}

func StringEncodedSenml(msg []byte) string {
	return fmt.Sprint(string(msg))
}
