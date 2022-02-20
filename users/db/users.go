package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/pseudoincorrect/bariot/users/models"
)

var _ models.UsersRepository = (*usersRepo)(nil) // static type check on userRepo
type usersRepo struct {
	db Database
}

func New(db *Database) models.UsersRepository {
	return &usersRepo{*db}
}

func (r *usersRepo) Save(ctx context.Context, t *models.User) (*models.User, error) {
	fail := func(err error) error {
		return fmt.Errorf("failed to save user: %v", err)
	}
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fail(err)
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
		return nil, fail(err)
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fail(err)
	}
	return t, nil
}

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

func (r *usersRepo) Delete(ctx context.Context, id string) (string, error) {
	fail := func(err error) error {
		return fmt.Errorf("failed to save user: %v", err)
	}
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", fail(err)
	}
	defer tx.Rollback(ctx)
	deletedId := ""
	// commandTag, err := r.db.conn.Exec(context.Background(),
	// 	"DELETE FROM users WHERE id=$1", id)
	err = tx.QueryRow(ctx, "DELETE FROM users WHERE id=$1 RETURNING id", id).Scan(&deletedId)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", fail(err)
	}
	return deletedId, nil
}

func (r *usersRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	fail := func(err error) error {
		return fmt.Errorf("failed to save user: %v", err)
	}
	tx, err := r.db.conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, fail(err)
	}
	defer tx.Rollback(ctx)
	var createdAt time.Time
	err = tx.QueryRow(ctx,
		"UPDATE users SET email=$1, full_name=$2 WHERE id=$3 RETURNING email, full_name, created_at",
		user.Email, user.FullName, user.Id).Scan(&user.Email, &user.FullName, &createdAt)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}
	if err = tx.Commit(ctx); err != nil {
		return nil, fail(err)
	}
	user.CreatedAt = createdAt.Format(time.RFC3339)
	return user, nil
}
