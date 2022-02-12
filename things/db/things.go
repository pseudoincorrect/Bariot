package db

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pseudoincorrect/bariot/things/models"
)

var _ models.ThingsRepository = (*thingsRepo)(nil) // static type check on thingRepo

type thingsRepo struct {
	db Database
}

func New(db *Database) models.ThingsRepository {
	return &thingsRepo{*db}
}

func (r *thingsRepo) Save(ctx context.Context, t *models.Thing) (*models.Thing, error) {
	fail := func(err error) error {
		return fmt.Errorf("failed to save thing: %v", err)
	}

	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return nil, fail(err)
	}
	defer tx.Rollback(ctx)

	// metadata_json, _ := json.Marshal(t.Metadata)

	res, err := tx.Exec(ctx, "INSERT INTO things ( key, name, user_id) VALUES ($1, $2, $3)",
		t.Key,
		t.Name,
		uuid.New())
	// t.Metadata

	fmt.Println(res)

	if err != nil {
		return nil, fail(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fail(err)
	}

	return t, nil
}

func (r *thingsRepo) Get(ctx context.Context, id string) (*models.Thing, error) {
	res := r.db.conn.QueryRow(ctx, "SELECT * from things WHERE id = $1", id)

	fmt.Println(res)

	return nil, nil
}

func (r *thingsRepo) Delete(ctx context.Context, id string) (*models.Thing, error) {
	return nil, nil
}

func (r *thingsRepo) Update(ctx context.Context, id string, thing *models.Thing) error {
	return nil
}
