package main

type Metadata map[string]interface{}

type ThingModel struct {
	id       string
	key      string
	name     string
	user_id  string
	metadata Metadata
}

type ThingRepository interface {
	saveThing(thing *ThingModel) error
	getThing(id string) (*ThingModel, error)
	deleteThing(id string) error
	updateThing(id string, thing *ThingModel) error
}
