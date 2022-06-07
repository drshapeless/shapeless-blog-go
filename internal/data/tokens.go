package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"errors"
	"time"

	"github.com/drshapeless/shapeless-blog/internal/validator"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
}

func generateToken(ttl time.Duration) (*Token, error) {
	token := &Token{
		Expiry: time.Now().Add(ttl),
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) Insert(token *Token) error {
	query := `
INSERT INTO tokens (hash, expiry)
VALUES (?, ?)`

	args := []interface{}{token.Hash, token.Expiry}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)
	return err

}

func (m TokenModel) New(ttl time.Duration) (*Token, error) {
	token, err := generateToken(ttl)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m TokenModel) GetByPlaintext(tokenPlaintext string) (*Token, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
SELECT hash, expiry
FROM tokens
WHERE hash = ? AND expiry > ?
`

	args := []interface{}{
		tokenHash[:],
		time.Now(),
	}

	var token Token

	err := m.DB.QueryRow(query, args...).Scan(
		&token.Hash,
		&token.Expiry,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &token, nil
}

func (m TokenModel) DeleteAllForExpired() error {
	query := `
DELETE FROM tokens
WHERE DATETIME('now', 'localtime') > expiry`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
