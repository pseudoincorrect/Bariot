package service

import (
	"context"
	"log"

	appErr "github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/things/models"
	"github.com/pseudoincorrect/bariot/things/rpc/client"
)

type ctxt context.Context

type Things interface {
	SaveThing(ctx ctxt, thingModel *models.Thing) (*models.Thing, error)
	GetThing(ctx ctxt, thingId string) (*models.Thing, error)
	DeleteThing(ctx ctxt, thingId string) (string, error)
	UpdateThing(ctx ctxt, thingModel *models.Thing) (*models.Thing, error)
	UserOfThingOrAdmin(ctx ctxt, token string, thingId string) (string, error)
	UserOfThingOnly(ctx ctxt, token string, thingId string) (string, error)
	UserOnly(ctx ctxt, token string) (string, error)
}

// type check on thingService
var _ Things = (*thingsService)(nil)

type thingsService struct {
	repository models.ThingsRepository
	auth       client.Auth
}

/// New creates a new thing service
func New(repository models.ThingsRepository, auth client.Auth) Things {
	return &thingsService{repository, auth}
}

/// SaveThing saves a thing to repository with thing model
func (s *thingsService) SaveThing(ctx ctxt, thing *models.Thing) (*models.Thing, error) {
	savedThing, err := s.repository.Save(ctx, thing)
	if err != nil {
		log.Println("Save Thing error:", err)
		return nil, err
	}
	return savedThing, nil
}

/// GetThing returns a thing from repository by id
func (s *thingsService) GetThing(ctx ctxt, id string) (*models.Thing, error) {
	thing, err := s.repository.Get(ctx, id)
	if err != nil {
		log.Println("Get Thing error:", err)
		return nil, err
	}
	return thing, nil
}

/// DeleteThing deletes a thing from repository by id
func (s *thingsService) DeleteThing(ctx ctxt, id string) (string, error) {
	resId, err := s.repository.Delete(ctx, id)
	if err != nil {
		log.Println("Delete Thing error:", err)
		return "", err
	}
	return resId, nil
}

/// UpdateThing updates a thing in repository by thing model
func (s *thingsService) UpdateThing(ctx ctxt, thing *models.Thing) (*models.Thing, error) {
	updatedThing, err := s.repository.Update(ctx, thing)
	if err != nil {
		log.Println("Update Thing error:", err)
		return nil, err
	}
	return updatedThing, nil
}

/// Check if the user is authorized to access the thing
/// return user ID , error if not a user or not the thing's user
func (s *thingsService) UserOfThingOnly(ctx ctxt, token string, thingId string) (string, error) {
	role, userId, err := s.auth.IsWhichUser(ctx, token)
	if err != nil {
		return "", appErr.ErrAuthentication
	}
	if role != "user" {
		return "", appErr.ErrAuthorization
	}
	thing, err := s.repository.Get(ctx, thingId)
	if err != nil {
		return "", appErr.ErrDb
	}
	if thing == nil {
		return "", appErr.ErrThingNotFound
	}
	if userId != thing.UserId {
		return "", appErr.ErrAuthorization
	}
	return userId, nil
}

/// Check if the user is authorized to access the thing
/// return user/admin ID, error if not a user or admin
func (s *thingsService) UserOfThingOrAdmin(ctx ctxt, token string, thingId string) (string, error) {
	role, userId, err := s.auth.IsWhichUser(ctx, token)
	if err != nil {
		return "", appErr.ErrAuthentication
	}
	if role == "admin" {
		return userId, nil
	}
	if role != "user" {
		return "", appErr.ErrAuthorization
	}
	thing, err := s.repository.Get(ctx, thingId)
	if err != nil {
		return "", appErr.ErrDb
	}
	if thing == nil {
		return "", appErr.ErrThingNotFound
	}
	if userId != thing.UserId {
		return "", appErr.ErrAuthorization
	}
	return userId, nil
}

/// Check if the token belong to a "user" user
/// return user id, error if not a user
func (s *thingsService) UserOnly(ctx ctxt, token string) (string, error) {
	role, userId, err := s.auth.IsWhichUser(ctx, token)
	if err != nil {
		return "", appErr.ErrAuthentication
	}
	if role != "user" {
		return "", appErr.ErrAuthorization
	}
	return userId, nil
}
