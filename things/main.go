package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	util "github.com/pseudoincorrect/bariot/things/utilities"
)

type config struct {
	mqttBrokerHost     string
	mqttBrokerHostHttp string
}

func loadConfig() config {
	var conf = config{
		mqttBrokerHost: util.GetEnv("MQTT_BROKER_HOST", "localhost:1883"),

		mqttBrokerHostHttp: util.GetEnv("MQTT_BROKER_HOST_HTTP", "localhost:8081"),
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

func main() {
	fmt.Println(" device service...")

	conf := loadConfig()

	client, err := connectMqttBroker(conf)
	
	if err != nil {
		fmt.Printf("connectMqttBroker error: %v", err)
		return
	}

	sub(client)
	publish(client)

	client.Disconnect(250)
}

func publish(client mqtt.Client) {
	num := 100
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(5 * time.Second)
	}
}

func sub(client mqtt.Client) {
	topic := "topic/test"
	token := client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}

func connectMqttBroker(conf config) (mqtt.Client, error) {
	var broker = "broker.emqx.io"
	var port = 1883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)

	tryCnt := 0

	for ! isMqttOnline(conf) {
		fmt.Println("Waiting for MQTT broker...")
		tryCnt++
		if tryCnt > 5 {
			return nil, errors.New("mqtt broker is not online")
		}
		time.Sleep(5 * time.Second)
	}

	token := client.Connect() 

	for token.Wait() && token.Error() != nil {
		return nil, token.Error();
	}

	// for token.Wait() && token.Error() != nil {
	// 	tryCnt++
	// 	if tryCnt > 5 {
	// 		return nil, token.Error()
	// 	}
	// 	time.Sleep(5 * time.Second)
	// 	token = client.Connect()
	// }
	return client, nil
}

func isMqttOnline(conf config) bool {
	res, err := http.Get(fmt.Sprintf("http://%s/status", conf.mqttBrokerHostHttp))
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
