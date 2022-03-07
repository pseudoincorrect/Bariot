package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/mainflux/senml"
)

func main() {
	// cfg := config{
	// 	mqttHost: "localhost",
	// 	mqttPort: "1883",
	// }
	// client, err := mqttConnect(cfg)
	createSenmlMsg()
}

type config struct {
	mqttHost string
	mqttPort string
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

func defaultMessageHandler(client mqtt.Client, msg mqtt.Message) {
	log.Printf("INCORRECT PUBLISH HERE: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
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
