package writer

import (
	"github.com/mainflux/senml"
)

type Writer interface {
	Write([]senml.Pack) error
}

type ThingData struct {
	ThingId     string
	SensorsData *senml.Pack
}
