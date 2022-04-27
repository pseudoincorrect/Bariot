package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go/v2"
)

type service struct {
	Writer Writer
	Mock   Mock
}

func main() {
	log.SetOutput(os.Stdout)
	s := service{}

	s.Writer = NewWriter()
	s.Mock = NewMock()

	err := s.Writer.Connect()
	if err != nil {
		log.Panic(err)
	}
	log.Println("Connected to InfluxDB")

	s.Mock.InitThings()

	go s.generateAndSend()

	// Wait for a CTRL-C to stop the program
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("Exiting Program...")
}

func (s *service) generateAndSend() {
	for {
		newStates := s.Mock.CreateArrayData()
		PrintStates(newStates)
		for _, state := range newStates {
			s.Writer.Write(&state)
		}
		time.Sleep(5 * time.Second)
	}
}

const INFLUX_ORG = "ggi_org"
const INFLUX_BUCKET = "ggi_bucket"
const INFLUX_TOKEN = "696A5C8CF1E5CBD65F480CF15773D1251FAC36FDCDA1D0119CC7DFC78DCCE064"
const INFLUX_HOST = "localhost"
const INFLUX_PORT = "8086"

type writer struct {
	influxClient influxdb.Client
}

type Writer interface {
	Connect() error
	Write(state *ThingState)
}

var _ Writer = (*writer)(nil)

func NewWriter() Writer {
	return &writer{}
}

// connectToInfluxdb setup a connection to an influxdb and check for health
func (w *writer) Connect() error {
	dbUrl := fmt.Sprintf("http://%s:%s", INFLUX_HOST, INFLUX_PORT)
	client := influxdb.NewClientWithOptions(dbUrl, INFLUX_TOKEN, influxdb.DefaultOptions().SetBatchSize(2))
	_, err := client.Health(context.Background())
	if err != nil {
		log.Panic("could not connect to influxdb")
	}
	w.influxClient = client
	return nil
}

// Write data (state) to influxdb
func (w *writer) Write(state *ThingState) {
	writeAPI := w.influxClient.WriteAPI(INFLUX_ORG, INFLUX_BUCKET)
	errChan := writeAPI.Errors()

	go func() {
		for err := range errChan {
			log.Println("Influxdb write error: ", err)
		}
	}()

	for _, r := range state.Measurements {
		p := influxdb.NewPointWithMeasurement(r.Name).
			AddTag("unit", r.Unit).
			AddTag("thingId", state.ThingId).
			SetTime(time.Unix(int64(state.Timestamp), 0))
		p.AddField("value", r.Value)

		writeAPI.WritePoint(p)
	}
	writeAPI.Flush()
}

const RANDOM_DATA_CNT = 10

type SensorData struct {
	Name  string
	Value int
	Unit  string
}

type ThingState struct {
	ThingId      string
	Timestamp    int64
	Measurements []SensorData
}

type mock struct {
	RandGen *rand.Rand
	States  []ThingState
}

type Mock interface {
	InitThings()
	CreateDataFromState(state ThingState, data *[]ThingState)
	CreateArrayData() []ThingState
	Print()
}

var _ Mock = (*mock)(nil)

func NewMock() Mock {
	return &mock{}
}

// InitThings init things' state and data
func (m *mock) InitThings() {

	ds := [...]string{
		"2cb1ca34-f912-4d8e-a2e1-8dd76bf9ce80",
		"22f2d1cd-be0d-4971-a1dd-260b27529c61",
		"4a2fd900-887b-4b17-a3f9-4b126179da17",
		"c1eed947-09f4-4792-9208-440c37120d0b",
		"df667ce7-87d4-49be-94ad-301a22956c8e",
		"e9115622-4a89-4ec4-9986-907c00ff393f",
		"e18b7944-e9eb-472e-aa14-922672958f43",
	}
	units := [][]string{
		{"temperature", "°C"},
		{"humidity", "%"},
		{"speed", "m/s"},
		{"failures", "f/h"},
		{"voltage", "V"},
	}

	for _, id := range ds {
		thingStates := ThingState{
			ThingId:      id,
			Timestamp:    time.Now().Unix(),
			Measurements: []SensorData{},
		}
		for _, unit := range units {
			sensorData := SensorData{
				Name:  unit[0],
				Unit:  unit[1],
				Value: 0,
			}
			thingStates.Measurements = append(thingStates.Measurements, sensorData)
		}
		m.States = append(m.States, thingStates)
	}
	randSource := rand.NewSource(time.Now().UnixNano())
	m.RandGen = rand.New(randSource)
}

// createDataFromState will update all measurements of a single thing from state
// and add these to a state array
func (m *mock) CreateDataFromState(currentState ThingState, states *[]ThingState) {
	for i, v := range currentState.Measurements {
		delta := m.RandGen.Intn(10) - 5
		newVal := v.Value + delta
		if newVal < -100 {
			newVal = -100
		}
		if newVal > 100 {
			newVal = 100
		}
		currentState.Measurements[i].Value = newVal
		currentState.Timestamp = time.Now().Unix()
	}
	*states = append(*states, currentState)
}

// createArrayData will pick randomly (can be twice the same thing),
// update randomly their state and create an array of ThingState
func (m *mock) CreateArrayData() []ThingState {
	states := []ThingState{}
	dataCnt := m.RandGen.Intn(7) + 3
	for i := 0; i < dataCnt; i++ {
		thingIndex := m.RandGen.Intn(len(m.States))
		m.CreateDataFromState(m.States[thingIndex], &states)
	}
	return states
}

// Print the things' states
func (m *mock) Print() {
	PrintStates(m.States)
}

// Print any thing's states
func PrintStates(states []ThingState) {
	log.Println()
	for _, state := range states {
		log.Println(state)
	}
	log.Println()
}
