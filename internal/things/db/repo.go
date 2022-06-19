package db

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pseudoincorrect/bariot/internal/things/models"
	e "github.com/pseudoincorrect/bariot/pkg/errors"
)

const uuidErr string = "invalid input syntax for type uuid"
const notFoundErr string = "no rows in result set"

// Static type checking
var _ models.ThingsRepository = (*thingsRepo)(nil)

type thingsRepo struct {
	db Database
}

func New(db *Database) models.ThingsRepository {
	return &thingsRepo{*db}
}

// Save a new thing to db
func (r *thingsRepo) Save(ctx context.Context, t *models.Thing) error {
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fail("save", err)
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
		return fail("save", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fail("save", err)
	}
	return nil
}

// Get a thing by id from db, update the thing given in the args
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
		return nil, e.ErrDbNotFound
	}
	if err != nil {
		return nil, e.ErrDb
	}
	thing.Id = thingUuid.String()
	thing.CreatedAt = createdAt.Format(time.RFC3339)
	return thing, nil
}

// Delete a thing by id from db
func (r *thingsRepo) Delete(ctx context.Context, id string) (string, error) {
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", fail("delete", err)
	}
	defer tx.Rollback(ctx)
	deletedId := ""

	err = tx.QueryRow(ctx, "DELETE FROM things WHERE id=$1 RETURNING id", id).Scan(&deletedId)

	if err != nil {
		log.Println("Error:", err)
		return "", fail("delete", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return "", fail("delete", err)
	}
	return deletedId, nil
}

// Get all things from db
func (r *thingsRepo) Update(ctx context.Context, thing *models.Thing) error {
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fail("update", err)
	}
	defer tx.Rollback(ctx)
	var createdAt time.Time

	err = tx.QueryRow(ctx,
		"UPDATE things SET key=$1, name=$2, user_id=$3 WHERE id=$4 RETURNING key, name, user_id, created_at",
		thing.Key, thing.Name, thing.UserId, thing.Id).Scan(&thing.Key, &thing.Name, &thing.UserId, &createdAt)

	if err != nil {
		log.Println("Error:", err)
		return fail("update", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fail("update", err)
	}
	// thing.CreatedAt = createdAt.Format(time.RFC3339)
	return nil
}

// Print and parse the DB error, return an app error
func fail(msg string, err error) error {
	// log.Println("DB failed", msg, err)
	if strings.Contains(err.Error(), notFoundErr) {
		return e.ErrDbNotFound
	}
	if strings.Contains(err.Error(), uuidErr) {
		return e.ErrDbUuid
	}
	return e.ErrDb
}
