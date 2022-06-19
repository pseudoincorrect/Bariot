package services

import (
	"encoding/json"
	"time"

	"github.com/pseudoincorrect/bariot/internal/reader/service"
	e "github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/stretchr/testify/mock"
)

type MockReader struct {
	mock.Mock
}

var _ service.Reader = (*MockReader)(nil)

func NewMockReader() MockReader {
	return MockReader{}
}

func (m *MockReader) AuthorizeSingleThing(userToken string, thingId string) error {
	args := m.Called()
	return args.Error(0)
}

type measurement struct {
	Timestamp string  `json:"timestamp"`
	Name      string  `json:"name"`
	Unit      string  `json:"unit"`
	Value     float64 `json:"value"`
	Tags      string  `json:"tags"`
}
type thingData struct {
	ThingId      string         `json:"thingId"`
	Measurements *[]measurement `json:"measurements"`
}

func (mr *MockReader) ReceiveThingData(thingId string, thingDataOut chan string, stop chan bool) {
	for {
		select {
		case <-stop:
			return
		default:
			seed := 0
			inc := 3
			measurements := createFakeMeasurements(seed, inc)
			seed += 3
			d := thingData{
				ThingId:      thingId,
				Measurements: measurements,
			}
			data, err := json.Marshal(&d)
			if err != nil {
				e.Handle(e.ErrParsing, err, "json measurements")
			}
			thingDataOut <- string(data)
		}
	}
}

func createFakeMeasurements(seed int, size int) *[]measurement {
	ms := make([]measurement, 0)
	for i := seed; i < size+seed; i++ {
		m := measurement{
			Timestamp: time.Now().String(),
			Name:      "temperature",
			Unit:      "celsius",
			Value:     float64(i),
			Tags:      "room 1",
		}
		ms = append(ms, m)
	}
	return &ms
}
