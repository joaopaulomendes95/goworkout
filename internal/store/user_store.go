package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// A password to be set into hash, or matched with a hash
type password struct {
	plainText *string
	hash      []byte
}

// Set hashes the given plaintext password using bcrypt and stores the hash.
// The plaintext password itself is also stored temporarily in the struct.
func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = hash
	return nil
}

// Matches compares a given plaintext password against the stored bcrypt hash.
// It returns true if they match, false otherwise.
func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

// User defines the structure for a user in the application.
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AnonymousUser represents an unauthenticated or guest user.
// This can be used in middleware to represent a user state before authentication.
var AnonymousUser = &User{}

// IsAnonymous checks if the user instance is the AnonymousUser.
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// PostgresUserStore implements the UserStore interface using a PostgreSQL database.
type PostgresUserStore struct {
	db *sql.DB
}

// NewPostgresUserStore creates a new instance of PostgresUserStore.
func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

// UserStore defines the interface for operations related to user data.
type UserStore interface {
	CreateUser(*User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
	GetUserToken(scope, tokenPlainText string) (*User, error)
}

// CreateUser inserts a new user record into the 'users' table.
// It hashes the user's password before storage and returns the user's ID and timestamps.
func (s *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username, email, password_hash, bio)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at
	`

	err := s.db.QueryRow(
		query,
		user.Username,
		user.Email,
		user.PasswordHash.hash,
		user.Bio,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByUsername retrieves a user from the database by their username.
// It returns the user details, including the password hash, or nil if not found.
func (s *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	query := `
	SELECT id, username, email, password_hash, bio, created_at, updated_at
	FROM users
	WHERE username = $1
	`

	err := s.db.QueryRow(
		query,
		username,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser modifies an existing user's details (username, bio) in the database.
// It updates the 'updated_at' timestamp.
func (s *PostgresUserStore) UpdateUser(user *User) error {
	query := `
	UPDATE users
	SET username = $1, bio = $2, updated_at = CURRENT_TIMESTAMP
	WHERE id = $3
	RETURNING updated_at
	`

	result, err := s.db.Exec(
		query,
		user.Username,
		user.Bio,
		user.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetUserToken retrieves a user based on a token's plaintext value and scope.
// The plaintext token is hashed (SHA256) before querying the database.
// This is used by the authentication middleware to validate tokens.
func (s *PostgresUserStore) GetUserToken(scope, plainTextPassword string) (*User, error) {
	// Hash the plaintext using SHA256
	tokenHash := sha256.Sum256([]byte(plainTextPassword))

	query := `
	SELECT u.id, u.username, u.email, u.password_hash, u.bio, u.created_at, u.updated_at
	FROM users u
	INNER JOIN tokens t ON t.user_id = u.id
	WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > $3
	`

	user := &User{
		PasswordHash: password{},
	}

	err := s.db.QueryRow(
		query,
		tokenHash[:],
		scope,
		time.Now(),
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}
