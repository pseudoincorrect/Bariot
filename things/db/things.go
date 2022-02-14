package db

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (r *thingsRepo) Delete(ctx context.Context, id string) (string, error) {
	fail := func(err error) error {
		return fmt.Errorf("failed to save thing: %v", err)
	}

	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return "", fail(err)
	}
	defer tx.Rollback(ctx)

	commandTag, err := r.db.conn.Exec(context.Background(),
		"DELETE FROM things WHERE id=$1", id)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	if commandTag.RowsAffected() != 1 {
		fmt.Println("Thing", id, "not found")
		return "", errors.New("thingNotFound")
	}

	if err = tx.Commit(ctx); err != nil {
		return "", fail(err)
	}

	return id, nil
}

func (r *thingsRepo) Update(ctx context.Context, thing *models.Thing) (*models.Thing, error) {
	fail := func(err error) error {
		return fmt.Errorf("failed to save thing: %v", err)
	}

	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})

	if err != nil {
		return nil, fail(err)
	}
	defer tx.Rollback(ctx)

	commandTag, err := r.db.conn.Exec(context.Background(),
		"UPDATE things SET key=$1, name=$2, user_id=$3 WHERE id=$4",
		thing.Key, thing.Name, thing.UserId, thing.Id)

	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	if commandTag.RowsAffected() != 1 {
		fmt.Println("Thing", thing.Id, "not found")
		return nil, errors.New("thingNotFound")
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fail(err)
	}

	return thing, nil
}
