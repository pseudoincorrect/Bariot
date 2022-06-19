// Service that make the link between MQTT (with auth) and NATS

package main

import (
	"log"
	"os"
	"time"

	"github.com/pseudoincorrect/bariot/internal/mqtt/mqtt"
	"github.com/pseudoincorrect/bariot/internal/mqtt/nats"
	auth "github.com/pseudoincorrect/bariot/pkg/auth/client"
	"github.com/pseudoincorrect/bariot/pkg/cache"
	"github.com/pseudoincorrect/bariot/pkg/env"
)

func main() {
	log.SetOutput(os.Stdout)
	conf := loadConfig()

	nats := nats.New(nats.Conf{
		Host: conf.natsHost,
		Port: conf.natsPort,
	})

	mqttSub := mqtt.New(mqtt.Conf{
		User:       conf.mqttUser,
		Pass:       conf.mqttPass,
		Host:       conf.mqttHost,
		Port:       conf.mqttPort,
		HealthPort: conf.mqttHealthPort})

	auth := auth.New(auth.Conf{
		Host: conf.rpcAuthHost,
		Port: conf.rpcAuthPort,
	})

	authCache := cache.New(cache.Conf{
		RedisHost: conf.redisHost,
		RedisPort: conf.redisPort,
	})

	err := authCache.Connect()
	if err != nil {
		log.Panic(err)
	}

	err = auth.StartAuthClient()
	if err != nil {
		log.Panic(err)
	}

	err = mqttSub.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer mqttSub.Disconnect()
	log.Printf("Connected to MQTT broker %s:%s\n", conf.mqttHost, conf.mqttPort)

	err = nats.Connect()
	if err != nil {
		log.Panic(err)
	}
	defer nats.Disconnect()

	const mqttThingsTopic = "things/#"
	const natsThingsSubject = "thingsMsg.>"
	natsHandler := nats.CreatePublisher(natsThingsSubject)
	mqttAuthorizer, err := mqtt.CreateAuthorizer(auth, authCache)
	if err != nil {
		log.Panic(err)
	}
	err = mqttSub.Subscriber(mqttThingsTopic, 0, mqttAuthorizer, natsHandler)
	if err != nil {
		log.Panic(err)
	}
	defer mqttSub.Disconnect()
	defer mqttSub.Unsubscribe(mqttThingsTopic)
	for {
		time.Sleep(5 * time.Second)
	}
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
	rpcAuthPort    string
	rpcAuthHost    string
	redisHost      string
	redisPort      string
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
		rpcAuthPort:    env.GetEnv("RPC_AUTH_PORT"),
		rpcAuthHost:    env.GetEnv("RPC_AUTH_HOST"),
		redisHost:      env.GetEnv("REDIS_HOST"),
		redisPort:      env.GetEnv("REDIS_PORT"),
	}
	return conf
}
