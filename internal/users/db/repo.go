package db

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pseudoincorrect/bariot/internal/users/models"
	e "github.com/pseudoincorrect/bariot/pkg/errors"
)

const uuidErr string = "invalid input syntax for type uuid"
const notFoundErr string = "no rows in result set"

// Static type checking
var _ models.UsersRepository = (*usersRepo)(nil)

type usersRepo struct {
	db Database
}

func New(db *Database) models.UsersRepository {
	return &usersRepo{*db}
}

// Save saves a user to the database
func (r *usersRepo) Save(ctx context.Context, t *models.User) error {
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fail("Save", err)
	}
	defer tx.Rollback(ctx)
	var id string
	var createdAt time.Time
	err = tx.QueryRow(ctx, "INSERT INTO users (email, full_name, hash_pass) VALUES ($1, $2, $3) RETURNING id, created_at ;",
		t.Email,
		t.FullName, t.HashPass).Scan(&id, &createdAt)
	t.Id = id
	t.CreatedAt = createdAt.Format(time.RFC3339)
	if err != nil {
		return fail("Save", err)
	}
	if err = tx.Commit(ctx); err != nil {
		return fail("Save", err)
	}
	return nil
}

// Get a user by id
func (r *usersRepo) Get(ctx context.Context, id string) (*models.User, error) {
	user := &models.User{}
	userUuid := uuid.UUID{}
	var createdAt time.Time
	row := r.db.conn.QueryRow(ctx, "SELECT * FROM users WHERE id::text=$1", id)
	err := row.Scan(
		&userUuid,
		&createdAt,
		&user.Email,
		&user.FullName,
		&user.HashPass,
		&user.Metadata,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.Id = userUuid.String()
	user.CreatedAt = createdAt.Format(time.RFC3339)
	return user, nil
}

// GetByEmail returns a user by email
func (r *usersRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
	userUuid := uuid.UUID{}
	var createdAt time.Time
	query := fmt.Sprintf("SELECT * FROM users WHERE email='%s';", email)
	row := r.db.conn.QueryRow(ctx, query)
	err := row.Scan(
		&userUuid,
		&createdAt,
		&user.Email,
		&user.FullName,
		&user.HashPass,
		&user.Metadata,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	user.Id = userUuid.String()
	user.CreatedAt = createdAt.Format(time.RFC3339)
	return user, nil
}

// Delete a user by id
func (r *usersRepo) Delete(ctx context.Context, id string) (string, error) {
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", fail("Delete", err)
	}
	defer tx.Rollback(ctx)
	deletedId := ""
	// commandTag, err := r.db.conn.Exec(context.Background(),
	// 	"DELETE FROM users WHERE id=$1", id)
	err = tx.QueryRow(ctx, "DELETE FROM users WHERE id=$1 RETURNING id", id).Scan(&deletedId)
	if err != nil {
		log.Println("Error:", err)
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", fail("Delete", err)
	}
	return deletedId, nil
}

// Update a user with a user model
func (r *usersRepo) Update(ctx context.Context, user *models.User) error {
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fail("Update", err)
	}
	defer tx.Rollback(ctx)
	var createdAt time.Time
	err = tx.QueryRow(ctx,
		"UPDATE users SET email=$1, full_name=$2 WHERE id=$3 RETURNING email, full_name, created_at",
		user.Email, user.FullName, user.Id).Scan(&user.Email, &user.FullName, &createdAt)
	if err != nil {
		log.Println("Error:", err)
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return fail("Update", err)
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
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
