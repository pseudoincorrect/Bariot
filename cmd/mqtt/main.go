package main

import (
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Println("Environment variable", key, "is not set")
		log.Println("Please set it and try again")
		panic("Environment variable " + key + " is not set and")
	}
	return value
}

type config struct {
	mqttHost  string
	mqttPort  string
	bariotEnv string
}

func loadConfig() config {
	var conf = config{
		mqttHost:  getEnv("MQTT_HOST"),
		mqttPort:  getEnv("MQTT_PORT"),
		bariotEnv: getEnv("BARIOT_ENV"),
	}
	return conf
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("TOPIC: %s\n", msg.Topic())
	log.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	config := loadConfig()
	time.Sleep(5 * time.Second)
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + config.mqttHost + ":" + config.mqttPort).SetClientID("emqx_test_client")

	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := c.Subscribe("testtopic/#", 0, nil); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}

	token := c.Publish("testtopic/1", 0, false, "Hello World")
	token.Wait()

	time.Sleep(6 * time.Second)

	if token := c.Unsubscribe("testtopic/#"); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		os.Exit(1)
	}

	c.Disconnect(250)
	time.Sleep(1 * time.Second)
}
