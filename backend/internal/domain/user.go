package domain

import (
	"crypto/sha256"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	PlainText *string
	Hash      []byte
}

func (p *Password) Set(plainTextPassword string) error {
	if plainTextPassword == "" {
		return errors.New("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}
	p.PlainText = &plainTextPassword
	p.Hash = hash
	return nil
}

func (p *Password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plainTextPassword))
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

func (p *Password) SetHash(hash []byte) {
	p.Hash = hash
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash Password  `json:"-"`
	Bio          string    `json:"bio,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == nil || u.ID == 0
}

type Token struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Hash      []byte    `json:"-"`
	Scope     string    `json:"scope"`
	Expiry    time.Time `json:"expiry"`
	CreatedAt time.Time `json:"created_at"`
}

func GenerateTokenHash(plainText string) []byte {
	hash := sha256.Sum256([]byte(plainText))
	return hash[:]
}

type UserWithToken struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}
