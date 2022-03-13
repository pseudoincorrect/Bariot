package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pseudoincorrect/bariot/pkg/auth/client/authClient"
	"github.com/pseudoincorrect/bariot/pkg/env"
	"github.com/pseudoincorrect/bariot/things/api"
	"github.com/pseudoincorrect/bariot/things/db"
	"github.com/pseudoincorrect/bariot/things/models"
	"github.com/pseudoincorrect/bariot/things/service"
)

type config struct {
	httpPort       string
	rpcAuthHost    string
	rpcAuthPort    string
	mqttHost       string
	mqttPort       string
	mqttStatusPort string
	dbHost         string
	dbPort         string
	dbUser         string
	dbPassword     string
	dbName         string
}

func loadConfig() config {
	var conf = config{
		httpPort:       env.GetEnv("HTTP_PORT"),
		rpcAuthHost:    env.GetEnv("RPC_AUTH_HOST"),
		rpcAuthPort:    env.GetEnv("RPC_AUTH_PORT"),
		mqttHost:       env.GetEnv("MQTT_HOST"),
		mqttPort:       env.GetEnv("MQTT_PORT"),
		mqttStatusPort: env.GetEnv("MQTT_STATUS_PORT"),
		dbHost:         env.GetEnv("PG_HOST"),
		dbPort:         env.GetEnv("PG_PORT"),
		dbName:         env.GetEnv("PG_DATABASE"),
		dbUser:         env.GetEnv("PG_USER"),
		dbPassword:     env.GetEnv("PG_PASSWORD"),
	}
	return conf
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connect lost: %v", err)
}

func publish(client mqtt.Client) {
	num := 100
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(30 * time.Second)
	}
}

func sub(client mqtt.Client) {
	topic := "topic/test"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	log.Printf("Subscribed to topic: %s", topic)
}

func connectMqttBroker(conf config) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", conf.mqttHost, conf.mqttPort))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)

	tryCnt := 0

	for !isMqttOnline(conf) {
		log.Println("Waiting for MQTT broker...")
		tryCnt++
		if tryCnt > 5 {
			return nil, errors.New("mqtt broker is not online")
		}
		time.Sleep(5 * time.Second)
	}

	token := client.Connect()

	for token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	return client, nil
}

func isMqttOnline(conf config) bool {
	res, err := http.Get(fmt.Sprintf("http://%s:%s/status", conf.mqttHost, conf.mqttStatusPort))
	if err != nil {
		log.Printf("isMqttOnline error: %v", err)
		return false
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("isMqttOnline error: %v", err)
		return false
	}
	log.Printf("isMqttOnline: %s", body)

	return true
}

func testDb() {
	log.Println("Test Db")
	conf := loadConfig()

	dbConf := db.DbConfig{
		Host:     conf.dbHost,
		Port:     conf.dbPort,
		Dbname:   conf.dbName,
		User:     conf.dbUser,
		Password: conf.dbPassword,
	}
	database, err := db.Init(dbConf)

	if err != nil {
		log.Println("Database Init error:", err)
	}
	thingsRepo := db.New(database)

	ctx := context.Background()

	thing := &models.Thing{
		Key:    "theKey",
		Name:   "thing1",
		UserId: uuid.New().String(),
	}

	thing, err = thingsRepo.Save(ctx, thing)

	if err != nil {
		log.Println("Save Thing error:", err)
		return
	}

	log.Println("Getting thing")
	thing, _ = thingsRepo.Get(ctx, thing.Id)
	log.Println(thing.String())

	log.Println("Updating thing")
	thing, _ = thingsRepo.Update(ctx, &models.Thing{
		Id:     thing.Id,
		Key:    "TheNewKey",
		Name:   "newThing1",
		UserId: uuid.New().String(),
	})

	thing, _ = thingsRepo.Get(ctx, thing.Id)

	log.Println("Deleting thing")
	thingsRepo.Delete(ctx, thing.Id)

	thing, _ = thingsRepo.Get(ctx, thing.Id)
	if thing == nil {
		log.Println("Success, no Thing found")
	}
}

func createService() (service.Things, error) {
	conf := loadConfig()

	dbConf := db.DbConfig{
		Host:     conf.dbHost,
		Port:     conf.dbPort,
		Dbname:   conf.dbName,
		User:     conf.dbUser,
		Password: conf.dbPassword,
	}
	database, err := db.Init(dbConf)

	if err != nil {
		log.Println("Database Init error:", err)
	}
	thingsRepo := db.New(database)

	authClientConf := authClient.AuthClientConf{
		Host: conf.rpcAuthHost,
		Port: conf.rpcAuthPort,
	}

	authClient := authClient.New(authClientConf)
	err = authClient.StartAuthClient()
	if err != nil {
		log.Println("Auth client error:", err)
		return nil, err
	}

	return service.New(thingsRepo, authClient), nil

}
func startHttp(s service.Things) error {
	conf := loadConfig()

	err := api.InitApi(conf.httpPort, s)
	return err
}

func main() {
	log.Println("Things service online")

	thingsService, err := createService()
	if err != nil {
		log.Panic("Create service error:", err)
	}

	err = startHttp(thingsService)
	if err != nil {
		log.Panic("Start http error:", err)
	}

	// client, err := connectMqttBroker(conf)
	// if err != nil {
	// 	log.Printf("connectMqttBroker error: %v", err)
	// 	return
	// }
	// sub(client)
	// publish(client)
	// client.Disconnect(250)
}
