package repos

import (
	"context"
	"errors"

	"github.com/pseudoincorrect/bariot/internal/things/models"
)

// Static type checking
var _ models.ThingsRepository = (*ThingsRepoMock)(nil)

type ThingsRepoMock struct {
	ThrowErr string
	Thing    *models.Thing
}

func NewRepoMock() ThingsRepoMock {
	return ThingsRepoMock{ThrowErr: ""}
}

// Save a new thing to db
func (r *ThingsRepoMock) Save(ctx context.Context, t *models.Thing) error {
	if r.ThrowErr != "" {
		return errors.New(r.ThrowErr)
	}
	return nil
}

// Get a thing by id from db
func (r *ThingsRepoMock) Get(ctx context.Context, id string) (*models.Thing, error) {
	if r.ThrowErr != "" {
		return nil, errors.New(r.ThrowErr)
	}
	return r.Thing, nil
}

// Delete a thing by id from db
func (r *ThingsRepoMock) Delete(ctx context.Context, id string) (string, error) {
	if r.ThrowErr != "" {
		return "", errors.New(r.ThrowErr)
	}
	return "", nil
}

// Get all things from db
func (r *ThingsRepoMock) Update(ctx context.Context, thing *models.Thing) error {
	if r.ThrowErr != "" {
		return errors.New(r.ThrowErr)
	}
	return nil
}
