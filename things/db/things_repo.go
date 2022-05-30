package db

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pseudoincorrect/bariot/things/models"
)

// Static type checking
var _ models.ThingsRepository = (*thingsRepo)(nil)

type thingsRepo struct {
	db Database
}

func New(db *Database) models.ThingsRepository {
	return &thingsRepo{*db}
}

// Save a new thing to db
func (r *thingsRepo) Save(ctx context.Context, t *models.Thing) (*models.Thing, error) {
	fail := func(err error) error {
		log.Println("failed to save thing:", err)
		return err
	}
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fail(err)
	}
	defer tx.Rollback(ctx)
	var id string
	var createdAt time.Time
	err = tx.QueryRow(ctx, "INSERT INTO things (key, name, user_id) VALUES ($1, $2, $3) returning id, created_at ;",
		t.Key,
		t.Name,
		t.UserId).Scan(&id, &createdAt)
	t.Id = id
	t.CreatedAt = createdAt.Format(time.RFC3339)
	if err != nil {
		return nil, fail(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fail(err)
	}
	return t, nil
}

// Get a thing by id from db
func (r *thingsRepo) Get(ctx context.Context, id string) (*models.Thing, error) {
	thing := &models.Thing{}
	thingUuid := uuid.UUID{}
	var createdAt time.Time
	row := r.db.conn.QueryRow(ctx, "SELECT * FROM things WHERE id::text=$1", id)
	err := row.Scan(
		&thingUuid,
		&createdAt,
		&thing.Key,
		&thing.Name,
		&thing.UserId,
		&thing.Metadata,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	thing.Id = thingUuid.String()
	thing.CreatedAt = createdAt.Format(time.RFC3339)
	return thing, nil
}

// Delete a thing by id from db
func (r *thingsRepo) Delete(ctx context.Context, id string) (string, error) {
	fail := func(err error) error {
		log.Printf("failed to save thing: %v", err)
		return err
	}
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", fail(err)
	}
	defer tx.Rollback(ctx)
	deletedId := ""
	err = tx.QueryRow(ctx, "DELETE FROM things WHERE id=$1 RETURNING id", id).Scan(&deletedId)
	if err != nil {
		log.Println("Error:", err)
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", fail(err)
	}
	return deletedId, nil
}

// Get all things from db
func (r *thingsRepo) Update(ctx context.Context, thing *models.Thing) (*models.Thing, error) {
	fail := func(err error) error {
		log.Printf("failed to save thing: %v", err)
		return err
	}
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fail(err)
	}
	defer tx.Rollback(ctx)
	var createdAt time.Time
	err = tx.QueryRow(ctx,
		"UPDATE things SET key=$1, name=$2, user_id=$3 WHERE id=$4 RETURNING key, name, user_id, created_at",
		thing.Key, thing.Name, thing.UserId, thing.Id).Scan(&thing.Key, &thing.Name, &thing.UserId, &createdAt)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fail(err)
	}
	thing.CreatedAt = createdAt.Format(time.RFC3339)
	return thing, nil
}
