package users

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, email, passwordHash string) error
	FindByEmail(ctx context.Context, email string) (id, hash string, err error)
	FindByID(ctx context.Context, id string) (User, error)

	// refresh tokens
	// StoreRefreshToken(ctx context.Context, userID, hash string) error
	// ValidateRefreshToken(hash string) (userID string, err error)
	// RevokeRefreshToken(ctx context.Context, hash string) error
}

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, email, passwordHash string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (email, password_hash) VALUES ($1, $2)`, email, passwordHash)
	return err
}

func (r *repo) FindByEmail(ctx context.Context, email string) (string, string, error) {
	var id, hash string
	err := r.db.QueryRow(ctx,
		`SELECT id, password_hash FROM users WHERE email=$1`,
		email,
	).Scan(&id, &hash)
	return id, hash, err
}

func (r *repo) FindByID(ctx context.Context, id string) (User, error) {
	var u User
	err := r.db.QueryRow(ctx,
		`SELECT id, email FROM users WHERE id=$1`,
		id,
	).Scan(&u.ID, &u.Email)
	return u, err
}

// func (r *repo) StoreRefreshToken(
// 	ctx context.Context,
// 	userID, hash string,
// ) error {
// 	_, err := r.db.Exec(ctx,
// 		`INSERT INTO refresh_tokens (user_id, token_hash)
// 		 VALUES ($1, $2)`,
// 		userID, hash,
// 	)
// 	return err
// }

// func (r *repo) ValidateRefreshToken(
// 	hash string,
// ) (string, error) {
// 	var userID string
// 	err := r.db.QueryRow(context.Background(),
// 		`SELECT user_id
// 		 FROM refresh_tokens
// 		 WHERE token_hash = $1`,
// 		hash,
// 	).Scan(&userID)

// 	return userID, err
// }

// func (r *repo) RevokeRefreshToken(
// 	ctx context.Context,
// 	hash string,
// ) error {
// 	_, err := r.db.Exec(ctx,
// 		`DELETE FROM refresh_tokens WHERE token_hash = $1`,
// 		hash,
// 	)
// 	return err
// }
