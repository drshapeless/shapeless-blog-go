package data

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
)

type Token struct {
	Plaintext string `json:"token"`
	Hash      []byte `json:"-"`
	UserID    int64  `json:"-"`
}

func generateToken(userID int64) (*Token, error) {
	token := &Token{
		UserID: userID,
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

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) New(userID int64) (*Token, error) {
	token, err := generateToken(userID)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	return token, err
}

func (m TokenModel) Insert(token *Token) error {
	query := `
INSERT INTO tokens (hash, user_id)
VALUES (?, ?)`

	_, err := m.DB.Exec(query, token.Hash, token.UserID)
	return err
}

func (m TokenModel) DeleteAllForUser(userID int64) error {
	query := `
DELETE FROM tokens
WHERE user_id = ?`

	_, err := m.DB.Exec(query, userID)
	return err
}
