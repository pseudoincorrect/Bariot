// Higher level function to manage things
// note that authorization is not checked here, it is checked in
// http hander. For instance we do not check that a user is authorized
// to update a thing here

package service

import (
	"context"

	"github.com/pseudoincorrect/bariot/internal/things/models"
	auth "github.com/pseudoincorrect/bariot/pkg/auth/client"
	rdb "github.com/pseudoincorrect/bariot/pkg/cache"
	e "github.com/pseudoincorrect/bariot/pkg/errors"
	"github.com/pseudoincorrect/bariot/pkg/utils/debug"
)

type Things interface {
	SaveThing(ctx context.Context, thingModel *models.Thing) (err error)
	GetThing(ctx context.Context, thingId string) (t *models.Thing, err error)
	DeleteThing(ctx context.Context, thingId string) (id string, err error)
	UpdateThing(ctx context.Context, thingModel *models.Thing) (err error)
	GetThingToken(ctx context.Context, thingId string, userId string) (token string, err error)
	UserOfThingOrAdmin(ctx context.Context, token string, thingId string) (userId string, err error)
	UserOfThingOnly(ctx context.Context, token string, thingId string) (userId string, err error)
	UserOnly(ctx context.Context, token string) (userId string, err error)
	GetUserOfThing(ctx context.Context, thingId string) (userId string, err error)
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
func (s *thingsService) SaveThing(
	ctx context.Context, thing *models.Thing) error {
	err := s.repository.Save(ctx, thing)
	if err != nil {
		debug.LogError("Save Thing error:", err)
		return err
	}
	return nil
}

// GetThing returns a thing from repository by id
func (s *thingsService) GetThing(ctx context.Context, id string) (*models.Thing, error) {
	thing, err := s.repository.Get(ctx, id)
	if err != nil {
		debug.LogError("Get Thing error:", err)
		return nil, err
	}
	return thing, nil
}

// DeleteThing deletes a thing from repository by id
func (s *thingsService) DeleteThing(ctx context.Context, id string) (string, error) {
	err := s.cache.DeleteTokenAndThingByThingId(id)
	if err != nil {
		debug.LogError("Could Delete and ThingId token in cache. err: ", err)
	}
	resId, err := s.repository.Delete(ctx, id)
	if err != nil {
		debug.LogError("Delete Thing error:", err)
		return "", err
	}
	return resId, nil
}

// UpdateThing updates a thing in repository by thing model
func (s *thingsService) UpdateThing(
	ctx context.Context, thing *models.Thing) error {
	err := s.repository.Update(ctx, thing)
	if err != nil {
		debug.LogError("Update Thing error:", err)
		return err
	}
	return nil
}

// GetThingToken return a JWT Token containing thing ID and user ID
func (s *thingsService) GetThingToken(
	ctx context.Context, thingId string, userId string) (string, error) {
	jwt, err := s.auth.GetThingToken(ctx, thingId, userId)
	if err != nil {
		debug.LogError("Get thing token error: ", err)
		return "", err
	}
	err = s.cache.SetTokenWithThingId(jwt, thingId)
	if err != nil {
		debug.LogError("Could not set token in cache. err: ", err)
	}
	return jwt, nil
}

// Check if the user is authorized to access the thing
// return user ID , error if not a user or not the thing's user
func (s *thingsService) UserOfThingOnly(
	ctx context.Context, token string, thingId string) (string, error) {
	role, userId, err := s.auth.IsWhichUser(ctx, token)
	if err != nil {
		return "", e.ErrAuthn
	}
	if role != "user" {
		return "", e.ErrAuthz
	}
	thing, err := s.repository.Get(ctx, thingId)
	if err != nil {
		return "", e.ErrDb
	}
	if thing == nil {
		return "", e.ErrNotFound
	}
	if userId != thing.UserId {
		return "", e.ErrAuthz
	}
	return userId, nil
}

// Check if the user is authorized to access the thing
// return user/admin ID, error if not a user or admin
func (s *thingsService) UserOfThingOrAdmin(
	ctx context.Context, token string, thingId string) (string, error) {
	role, userId, err := s.auth.IsWhichUser(ctx, token)
	if err != nil {
		return "", e.ErrAuthn
	}
	if role == "admin" {
		return userId, nil
	}
	if role != "user" {
		return "", e.ErrAuthz
	}
	thing, err := s.repository.Get(ctx, thingId)
	if err != nil {
		return "", e.ErrDb
	}
	if thing == nil {
		return "", e.ErrNotFound
	}
	if userId != thing.UserId {
		return "", e.ErrAuthz
	}
	return userId, nil
}

// Check if the token belong to a "user" user
// return user id, error if not a user
func (s *thingsService) UserOnly(ctx context.Context, token string) (string, error) {
	role, userId, err := s.auth.IsWhichUser(ctx, token)
	if err != nil {
		return "", e.ErrAuthn
	}
	if role != "user" {
		return "", e.ErrAuthz
	}
	return userId, nil
}

// Return the User ID of a given Thing ID
func (s *thingsService) GetUserOfThing(ctx context.Context, thingId string) (string, error) {
	thing, err := s.GetThing(ctx, thingId)
	if err != nil {
		return "", err
	}
	return thing.UserId, nil
}
