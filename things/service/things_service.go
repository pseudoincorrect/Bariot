// Higher level function to manage things
// note that authorization is not checked here, it is checked in
// http hander. For instance we do not check that a user is authorized
// to update a thing here

package service

import (
	"context"
	"log"

	auth "github.com/pseudoincorrect/bariot/pkg/auth/client"
	rdb "github.com/pseudoincorrect/bariot/pkg/cache"
	appErr "github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/things/models"
)

type ctxt context.Context

type Things interface {
	SaveThing(ctx ctxt, thingModel *models.Thing) (*models.Thing, error)
	GetThing(ctx ctxt, thingId string) (*models.Thing, error)
	DeleteThing(ctx ctxt, thingId string) (string, error)
	UpdateThing(ctx ctxt, thingModel *models.Thing) (*models.Thing, error)
	GetThingToken(ctx ctxt, thingId string, userId string) (string, error)
	UserOfThingOrAdmin(ctx ctxt, token string, thingId string) (string, error)
	UserOfThingOnly(ctx ctxt, token string, thingId string) (string, error)
	UserOnly(ctx ctxt, token string) (string, error)
}

// type check on thingService
var _ Things = (*thingsService)(nil)

type thingsService struct {
	repository models.ThingsRepository
	auth       auth.Auth
	cache      rdb.ThingCache
}

// New creates a new thing service
func New(repository models.ThingsRepository, auth auth.Auth, cache rdb.ThingCache) Things {
	return &thingsService{repository, auth, cache}
}

// SaveThing saves a thing to repository with thing model
func (s *thingsService) SaveThing(ctx ctxt, thing *models.Thing) (*models.Thing, error) {
	savedThing, err := s.repository.Save(ctx, thing)
	if err != nil {
		log.Println("Save Thing error:", err)
		return nil, err
	}
	return savedThing, nil
}

// GetThing returns a thing from repository by id
func (s *thingsService) GetThing(ctx ctxt, id string) (*models.Thing, error) {
	thing, err := s.repository.Get(ctx, id)
	if err != nil {
		log.Println("Get Thing error:", err)
		return nil, err
	}
	return thing, nil
}

// DeleteThing deletes a thing from repository by id
func (s *thingsService) DeleteThing(ctx ctxt, id string) (string, error) {
	err := s.cache.DeleteTokenAndTokenByThingId(id)
	if err != nil {
		log.Println("Could Delete and ThingId token in cache. err: ", err)
	}
	resId, err := s.repository.Delete(ctx, id)
	if err != nil {
		log.Println("Delete Thing error:", err)
		return "", err
	}
	return resId, nil
}

// UpdateThing updates a thing in repository by thing model
func (s *thingsService) UpdateThing(ctx ctxt, thing *models.Thing) (*models.Thing, error) {
	updatedThing, err := s.repository.Update(ctx, thing)
	if err != nil {
		log.Println("Update Thing error:", err)
		return nil, err
	}
	return updatedThing, nil
}

// GetThingToken return a JWT Token containing thing ID and user ID
func (s *thingsService) GetThingToken(ctx ctxt, thingId string, userId string) (string, error) {
	jwt, err := s.auth.GetThingToken(ctx, thingId, userId)
	if err != nil {
		log.Println("Get thing token error: ", err)
		return "", err
	}
	err = s.cache.SetTokenWithThingId(jwt, thingId)
	if err != nil {
		log.Println("Could not set token in cache. err: ", err)
	}
	return jwt, nil
}

// Check if the user is authorized to access the thing
// return user ID , error if not a user or not the thing's user
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

// Check if the user is authorized to access the thing
// return user/admin ID, error if not a user or admin
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

// Check if the token belong to a "user" user
// return user id, error if not a user
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
