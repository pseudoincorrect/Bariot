package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pseudoincorrect/bariot/pkg/utils"
)

func TestSenmlMsgFormt(t *testing.T) {
	msg, err := createSenmlMsg()
	if err != nil {
		t.Fatal("error creating Senml message", err)
	}
	if !json.Valid(msg) {
		t.Fatal("\nInvalid Senml msg\n", string(msg))
	}
}

func TestSendSenmlOverMqtt(t *testing.T) {
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
