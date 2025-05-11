package store

import (
	"database/sql"
	"time"

	"github.com/strangecousinwst/goworkout/internal/tokens"
)

// PostgresTokenStore implements the TokenStore interface using a PostgreSQL database.
type PostgresTokenStore struct {
	db *sql.DB
}

// NewPostgresTokenStore creates a new instance of PostgresTokenStore.
func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

// TokenStore defines the interface for operations related to authentication tokens.
type TokenStore interface {
	// Insert stores a new token in the database.
	Insert(token *tokens.Token) error
	// CreateNewToken generates a new token for a given user ID, then inserts it into the database.
	CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error)
	// DeleteAllTokensForUser removes all tokens of a specific scope for a given user.
	DeleteAllTokensForUser(userID int, scope string) error
}

// CreateNewToken generates a new token, hashes it, and inserts it into the database.
// It returns the generated token (including its plaintext version) or an error.
func (t *PostgresTokenStore) CreateNewToken(userID int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	err = t.Insert(token)
	if err != nil {
		return nil, nil
	}

	return token, nil
}

// Insert adds a new token record to the 'tokens' table.
// The token's hash, user ID, expiry time, and scope are stored.
func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)
	`

	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

// DeleteAllTokensForUser removes all tokens associated with a specific user ID and scope.
// For example, this can invalidate all 'authentication' tokens for a user.
func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `
	DELETE FROM tokens
	WHERE scope = $1 AND user_id = $2
	`

	_, err := t.db.Exec(query, scope, userID)
	return err
}
