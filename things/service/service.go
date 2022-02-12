package service

import (
	"github.com/pseudoincorrect/bariot/things/models"
)

type Service interface {
	SaveThing(thing *models.Thing) error
	GetThing(id string) (*models.Thing, error)
	DeleteThing(id string) error
	UpdateThing(id string, thing *models.Thing) error
}

var _ Service = (*thingsService)(nil) // type check on thingService

type thingsService struct {
	repository models.ThingsRepository
}

func (s *thingsService) SaveThing(thing *models.Thing) error {
	return nil
}

func (s *thingsService) GetThing(id string) (*models.Thing, error) {
	return nil, nil
}

func (s *thingsService) DeleteThing(id string) error {
	return nil
}

func (s *thingsService) UpdateThing(id string, thing *models.Thing) error {
	return nil
}
