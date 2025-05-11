package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

// ScopeAuth defines the scope for authentication tokens.
// Scopes can be used to differentiate tokens for various purposes (e.g., auth, password reset).
const (
	ScopeAuth = "authentication"
)

// Token represents an authentication token, including its plaintext value (for sending to the client once),
// its SHA256 hash (for storing in the database), user association, expiry, and scope.
type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int       `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// GenerateToken creates a new unique token for a given user ID, time-to-live (ttl), and scope.
// It generates a cryptographically secure random string for the token's plaintext value
// and a SHA256 hash of this plaintext for database storage.
func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Generate 32 random bytes. This will result in a base32 encoded string of
	// roughly 52 characters, which is a good length for a token.
	emptyBytes := make([]byte, 32)
	_, err := rand.Read(emptyBytes)
	if err != nil {
		// This error is serious as it means the crypto/rand source failed.
		return nil, err
	}

	// Encode the random bytes to a base32 string. Base32 is URL-safe and case-insensitive.
	// NoPadding is used to make the token slightly shorter.
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyBytes)

	// Generate a SHA256 hash of the plaintext token. This hash is stored in the database.
	// Storing the hash instead of the plaintext token enhances security: if the database
	// is compromised, the actual tokens are not exposed.
	hash := sha256.Sum256([]byte(token.PlainText))

	// Convert the [32]byte array to a []byte slice.
	token.Hash = hash[:]

	return token, nil
}
