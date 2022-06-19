package services

import (
	c "context"
	"errors"

	"github.com/pseudoincorrect/bariot/internal/things/models"
	"github.com/pseudoincorrect/bariot/internal/things/service"
)

type ThingsMock struct {
	Things   []*models.Thing
	Thing    *models.Thing
	Token    string
	ThrowErr string
	UserId   string
	ThingId  string
}

var _ service.Things = (*ThingsMock)(nil)

func NewThingsMock() ThingsMock {
	return ThingsMock{}
}

func (s *ThingsMock) SaveThing(ctx c.Context, thingModel *models.Thing) error {
	if s.ThrowErr != "" {
		return errors.New(s.ThrowErr)
	}
	s.Things = append(s.Things, thingModel)
	return nil
}

func (s *ThingsMock) GetThing(ctx c.Context, thingId string) (*models.Thing, error) {
	if s.ThrowErr != "" {
		return nil, errors.New(s.ThrowErr)
	}
	if s.Thing != nil && thingId == s.Thing.Id {
		return s.Thing, nil
	}
	return nil, nil
}

func (s *ThingsMock) DeleteThing(ctx c.Context, thingId string) (string, error) {
	if s.ThrowErr != "" {
		return "", errors.New(s.ThrowErr)
	}
	return thingId, nil
}

func (s *ThingsMock) UpdateThing(ctx c.Context, thingModel *models.Thing) error {
	return nil
}

func (s *ThingsMock) GetThingToken(ctx c.Context, thingId string, userId string) (string, error) {
	if s.ThrowErr != "" {
		return "", errors.New(s.ThrowErr)
	}
	return s.Token, nil
}

func (s *ThingsMock) UserOfThingOrAdmin(ctx c.Context, token string, thingId string) (string, error) {
	if s.ThrowErr != "" {
		return "", errors.New(s.ThrowErr)
	}
	return s.UserId, nil
}

func (s *ThingsMock) UserOfThingOnly(ctx c.Context, token string, thingId string) (string, error) {
	if s.ThrowErr != "" {
		return "", errors.New(s.ThrowErr)
	}
	return s.UserId, nil
}

func (s *ThingsMock) UserOnly(ctx c.Context, token string) (string, error) {
	if s.ThrowErr != "" {
		return "", errors.New(s.ThrowErr)
	}
	return s.UserId, nil
}

func (s *ThingsMock) GetUserOfThing(ctx c.Context, userId string) (string, error) {
	if s.ThrowErr != "" {
		return "", errors.New(s.ThrowErr)
	}
	return s.UserId, nil
}
