package main

import (
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pseudoincorrect/bariot/pkg/env"
)

const thingsTopic = "things/#"

type config struct {
	mqttHost  string
	mqttPort  string
	bariotEnv string
}

func loadConfig() config {
	var conf = config{
		mqttHost:  env.GetEnv("MQTT_HOST"),
		mqttPort:  env.GetEnv("MQTT_PORT"),
		bariotEnv: env.GetEnv("BARIOT_ENV"),
	}
	return conf
}

var defaultMessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("INCORRECT PUBLISH HERE: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
}

func thingsHandler(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("[%s]  ", msg.Topic())
	fmt.Printf("%s\n", msg.Payload())
}

func main() {
	config := loadConfig()
	time.Sleep(5 * time.Second)
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + config.mqttHost + ":" + config.mqttPort).SetClientID("emqx_client_1")

	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(defaultMessageHandler)
	opts.SetPingTimeout(1 * time.Minute)

	c := mqtt.NewClient(opts)
	token := c.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	token = c.Subscribe(thingsTopic, 0, thingsHandler)
	if token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}

	// token = c.Publish("testtopic/1", 0, false, "Hello World")
	// token.Wait()

	time.Sleep(60 * time.Minute)

	token = c.Unsubscribe(thingsTopic)
	if token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
	time.Sleep(1 * time.Second)
}
