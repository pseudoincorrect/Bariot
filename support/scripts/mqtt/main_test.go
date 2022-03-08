package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/pseudoincorrect/bariot/pkg/utils"
)

// Publish Mqtt message, test without cache
// go test -run TestSendSenmlOverMqtt github.com/pseudoincorrect/bariot/support/scripts/mqtt -v -count=1

func TestCreateSenmlMsg(t *testing.T) {
	msg, err := createSenmlMsg()
	if err != nil {
		t.Fatal("error creating Senml message", err)
	}
	if !json.Valid(msg) {
		t.Fatal("\nInvalid Senml msg\n", string(msg))
	}
}

func TestSendSenmlOverMqtt(t *testing.T) {
	log.SetOutput(os.Stdout)
	err := MqttConnectAndSend()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMqttHealthCheck(t *testing.T) {
	resp, err := http.Get("http://admin:public@localhost:8084/api/v4/brokers")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatal("could not connect")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(utils.PrettyJsonString(string(body)))
}
