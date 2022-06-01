package mock

import (
	c "context"
	"errors"

	"github.com/pseudoincorrect/bariot/things/models"
	"github.com/pseudoincorrect/bariot/things/service"
)

type thingsMock struct {
	Things   []*models.Thing
	Thing    *models.Thing
	ThrowErr string
}

var _ service.Things = (*thingsMock)(nil)

func New() thingsMock {
	return thingsMock{}
}

func (s *thingsMock) SaveThing(ctx c.Context, thingModel *models.Thing) (*models.Thing, error) {
	s.Things = append(s.Things, thingModel)
	return thingModel, nil
}

func (s *thingsMock) GetThing(ctx c.Context, thingId string) (*models.Thing, error) {
	if s.ThrowErr != "" {
		return nil, errors.New(s.ThrowErr)
	}
	if s.Thing != nil && thingId == s.Thing.Id {
		return s.Thing, nil
	}
	return nil, nil
}

func (s *thingsMock) DeleteThing(ctx c.Context, thingId string) (string, error) {
	return "", nil
}

func (s *thingsMock) UpdateThing(ctx c.Context, thingModel *models.Thing) (*models.Thing, error) {
	return nil, nil
}

func (s *thingsMock) GetThingToken(ctx c.Context, thingId string, userId string) (string, error) {
	return "", nil
}

func (s *thingsMock) UserOfThingOrAdmin(ctx c.Context, token string, thingId string) (string, error) {
	return "", nil
}

func (s *thingsMock) UserOfThingOnly(ctx c.Context, token string, thingId string) (string, error) {
	return "", nil
}

func (s *thingsMock) UserOnly(ctx c.Context, token string) (string, error) {
	return "", nil
}

func CreateMockThing() models.Thing {
	return models.Thing{
		Id:        "00000000-0000-0000-0000-000000000001",
		CreatedAt: "2022-06-01T14:35:40+03:00",
		Key:       "123456789",
		Name:      "Thing_1",
		UserId:    "00000000-0000-0000-0000-000000000002",
		Metadata:  models.Metadata{"unit": "temperature"},
	}
}
