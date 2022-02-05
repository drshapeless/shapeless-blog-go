package data

import (
	"database/sql"
	"errors"
)

type Tag struct {
	Name    string `json:"name"`
	PostID  string `json:"post_id"`
	Version int64  `json:"-"`
}

type TagModel struct {
	DB *sql.DB
}

func (m *TagModel) Insert(t *Tag) error {
	query := `
INSERT INTO tags(name, post_id)
VALUES (?, ?)
RETURNING version`

	return m.DB.QueryRow(query, t.Name, t.PostID).Scan(&t.Version)
}

func (m *TagModel) Get(name string) (*Tag, error) {
	if name == "" {
		return nil, ErrRecordNotFound
	}

	query := `
SELECT name, post_id, version
FROM tags
WHERE name = ?`

	var t Tag

	err := m.DB.QueryRow(query, name).Scan(
		&t.Name,
		&t.PostID,
		&t.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &t, nil
}

func (m *TagModel) Update(t *Tag) error {
	query := `
UPDATE tags
SET post_id = ?, version = version + 1
WHERE name = ? AND version = ?`

	args := []interface{}{
		t.PostID,
		t.Name,
		t.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&t.Version)
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

func (m *TagModel) Delete(name string) error {
	if name == "" {
		return ErrRecordNotFound
	}

	query := `
DELETE FROM tags
WHERE name = ?`

	result, err := m.DB.Exec(query, name)
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
