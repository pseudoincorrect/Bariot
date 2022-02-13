package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pseudoincorrect/bariot/things/api"
	"github.com/pseudoincorrect/bariot/things/db"
	"github.com/pseudoincorrect/bariot/things/models"
	"github.com/pseudoincorrect/bariot/things/service"
	util "github.com/pseudoincorrect/bariot/things/utilities"
)

type config struct {
	httpPort       string
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
		httpPort:       util.GetEnv("HTTP_PORT"),
		mqttHost:       util.GetEnv("MQTT_HOST"),
		mqttPort:       util.GetEnv("MQTT_PORT"),
		mqttStatusPort: util.GetEnv("MQTT_STATUS_PORT"),
		dbHost:         util.GetEnv("PG_HOST"),
		dbPort:         util.GetEnv("PG_PORT"),
		dbName:         util.GetEnv("PG_DATABASE"),
		dbUser:         util.GetEnv("PG_USER"),
		dbPassword:     util.GetEnv("PG_PASSWORD"),
	}
	return conf
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
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
	fmt.Printf("Subscribed to topic: %s", topic)
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
		fmt.Println("Waiting for MQTT broker...")
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
		fmt.Printf("isMqttOnline error: %v", err)
		return false
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("isMqttOnline error: %v", err)
		return false
	}
	fmt.Printf("isMqttOnline: %s", body)

	return true
}

func testDb() {
	fmt.Println("Test Db")
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
		fmt.Println("Database Init error:", err)
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
		fmt.Println("Save Thing error:", err)
		return
	}

	fmt.Println("Getting thing")
	thing, _ = thingsRepo.Get(ctx, thing.Id)
	fmt.Println(thing.String())

	fmt.Println("Updating thing")
	thing, _ = thingsRepo.Update(ctx, &models.Thing{
		Id:     thing.Id,
		Key:    "TheNewKey",
		Name:   "newThing1",
		UserId: uuid.New().String(),
	})

	thing, _ = thingsRepo.Get(ctx, thing.Id)

	fmt.Println("Deleting thing")
	thingsRepo.Delete(ctx, thing.Id)

	thing, _ = thingsRepo.Get(ctx, thing.Id)
	if thing == nil {
		fmt.Println("Success, no Thing found")
	}
}

func createService() service.Things {
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
		fmt.Println("Database Init error:", err)
	}
	thingsRepo := db.New(database)

	return service.New(thingsRepo)

}
func testHttp(s service.Things) {
	conf := loadConfig()

	api.InitApi(conf.httpPort, s)
}

func main() {
	fmt.Println(" device service...")

	thingsService := createService()

	testHttp(thingsService)

	// client, err := connectMqttBroker(conf)
	// if err != nil {
	// 	fmt.Printf("connectMqttBroker error: %v", err)
	// 	return
	// }
	// sub(client)
	// publish(client)
	// client.Disconnect(250)
}
