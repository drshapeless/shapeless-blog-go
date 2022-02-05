package data

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string   `json:"username"`
	Password password `json:"-"`
	Version  int64    `json:"-"`
}

type password struct {
	plaintext *string
	hash      []byte
}

func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, nil
		}
	}

	return true, nil
}

// This is a good to have feature, but very useless.
// There is no point in creating users for a personal blog.
// I am the only user.
type UserModel struct {
	DB *sql.DB
}

func (m UserModel) Insert(u *User) error {
	query := `
INSERT INTO users (username, password_hash)
VALUES (?, ?)
RETURNING version`

	err := m.DB.QueryRow(query, u.Username, u.Password.hash).Scan(&u.Version)
	if err != nil {
		// Should check duplicated username here!
		return err
	}

	return nil
}

func (m UserModel) Get(name string) (*User, error) {
	query := `
SELECT username, password_hash, version
FROM users
WHERE username = ?`

	var u User
	err := m.DB.QueryRow(query, name).Scan(&u.Username, &u.Password.hash, &u.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &u, nil
}

func (m UserModel) Update(u *User) error {
	query := `
UPDATE users
SET password_hash = ?, version = version + 1
WHERE username = ? AND version = ?
RETURNING version`

	err := m.DB.QueryRow(query, u.Password.hash, u.Username, u.Version).Scan(&u.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
