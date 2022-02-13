package service

import (
	"context"
	"fmt"

	"github.com/pseudoincorrect/bariot/things/models"
)

type Things interface {
	SaveThing(ctx context.Context, thing *models.Thing) (*models.Thing, error)
	GetThing(ctx context.Context, id string) (*models.Thing, error)
	DeleteThing(ctx context.Context, id string) (string, error)
	UpdateThing(ctx context.Context, thing *models.Thing) (*models.Thing, error)
}

// type check on thingService
var _ Things = (*thingsService)(nil)

type thingsService struct {
	repository models.ThingsRepository
}

/// New creates a new thing service
func New(repository models.ThingsRepository) Things {
	return &thingsService{repository}
}

/// SaveThing saves a thing to repository with thing model
func (s *thingsService) SaveThing(ctx context.Context, thing *models.Thing) (*models.Thing, error) {
	savedThing, err := s.repository.Save(ctx, thing)
	if err != nil {
		fmt.Println("Save Thing error:", err)
		return nil, err
	}
	return savedThing, nil
}

/// GetThing returns a thing from repository by id
func (s *thingsService) GetThing(ctx context.Context, id string) (*models.Thing, error) {
	thing, err := s.repository.Get(ctx, id)
	if err != nil {
		fmt.Println("Get Thing error:", err)
		return nil, err
	}
	return thing, nil
}

/// DeleteThing deletes a thing from repository by id
func (s *thingsService) DeleteThing(ctx context.Context, id string) (string, error) {
	resId, err := s.repository.Delete(ctx, id)
	if err != nil {
		fmt.Println("Delete Thing error:", err)
		return "", err
	}
	return resId, nil
}

/// UpdateThing updates a thing in repository by thing model
func (s *thingsService) UpdateThing(ctx context.Context, thing *models.Thing) (*models.Thing, error) {
	updatedThing, err := s.repository.Update(ctx, thing)
	if err != nil {
		fmt.Println("Update Thing error:", err)
		return nil, err
	}
	return updatedThing, nil
}
