package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joaodias/golang-svelte-fullstack/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
}

type TokenRepository interface {
	Create(ctx context.Context, token *domain.Token) error
	GetByHashAndScope(ctx context.Context, hash []byte, scope string) (*domain.Token, error)
	DeleteByUserID(ctx context.Context, userID int) error
}

type PostgresUserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, bio)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.PasswordHash.Hash,
		user.Bio,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *PostgresUserRepo) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := &domain.User{PasswordHash: domain.Password{}}
	query := `SELECT id, username, email, password_hash, bio, created_at, updated_at FROM users WHERE username = $1`
	err := r.db.QueryRow(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.Hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *PostgresUserRepo) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user := &domain.User{PasswordHash: domain.Password{}}
	query := `SELECT id, username, email, password_hash, bio, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.Hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

type PostgresTokenRepo struct {
	db *pgxpool.Pool
}

func NewTokenRepository(db *pgxpool.Pool) *PostgresTokenRepo {
	return &PostgresTokenRepo{db: db}
}

func (r *PostgresTokenRepo) Create(ctx context.Context, token *domain.Token) error {
	query := `
		INSERT INTO tokens (user_id, hash, scope, expiry)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query,
		token.UserID,
		token.Hash,
		token.Scope,
		token.Expiry,
	).Scan(&token.ID, &token.CreatedAt)
}

func (r *PostgresTokenRepo) GetByHashAndScope(ctx context.Context, hash []byte, scope string) (*domain.Token, error) {
	token := &domain.Token{}
	query := `SELECT id, user_id, hash, scope, expiry, created_at FROM tokens WHERE hash = $1 AND scope = $2 AND expiry > $3`
	err := r.db.QueryRow(ctx, query, hash, scope, time.Now()).Scan(
		&token.ID,
		&token.UserID,
		&token.Hash,
		&token.Scope,
		&token.Expiry,
		&token.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (r *PostgresTokenRepo) DeleteByUserID(ctx context.Context, userID int) error {
	_, err := r.db.Exec(ctx, "DELETE FROM tokens WHERE user_id = $1", userID)
	return err
}
