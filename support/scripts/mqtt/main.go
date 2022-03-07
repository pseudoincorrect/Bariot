package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mainflux/senml"
	"github.com/pseudoincorrect/bariot/pkg/errors"
)

type mqttTester struct {
	conf   config
	client mqtt.Client
}

func main() {
	MqttConnectAndSend()
}

func MqttConnectAndSend() error {
	var m mqttTester
	m.conf = config{
		mqttHost:       "localhost",
		mqttPort:       "1883",
		mqttUser:       "admin",
		mqttPass:       "public",
		mqttHealthPort: "8084",
	}
	err := m.mqttConnect()
	if err != nil {
		log.Panic("could not connect to MQTT broker")
		// return err
	}
	defer m.mqttDisconnect()
	msg, _ := createSenmlMsg()
	err = m.mqttPublish(string(msg))
	if err != nil {
		log.Panic("could not publish MQTT message")
		// return err
	}
	time.Sleep(1 * time.Second)
	return nil
}

type config struct {
	mqttHost       string
	mqttPort       string
	mqttUser       string
	mqttPass       string
	mqttHealthPort string
}

func (m *mqttTester) mqttSetOpts() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + m.conf.mqttHost + ":" + m.conf.mqttPort).SetClientID("bariot_mqtt_things")
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetPingTimeout(1 * time.Minute)
	return opts
}

func (m *mqttTester) mqttHealthCheckBlocking() {
	for {
		err := m.mqttHealthCheck()
		if err == nil {
			break
		}
		fmt.Println("MQTT broker not online, retrying later...")
		time.Sleep(5 * time.Second)
	}
}

func (m *mqttTester) mqttHealthCheck() error {
	url := "http://" + m.conf.mqttUser + ":" + m.conf.mqttPass + "@" +
		m.conf.mqttHost + ":" + m.conf.mqttHealthPort + "/api/v4/brokers"
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

func (m *mqttTester) mqttConnect() error {
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := m.mqttSetOpts()
	m.mqttHealthCheckBlocking()
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

func (m *mqttTester) mqttPublish(msg string) error {
	topic := "things/123456789"
	token := m.client.Publish(topic, 0, false, msg)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
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
