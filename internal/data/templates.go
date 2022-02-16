package data

import (
	"database/sql"
	"errors"
)

type Template struct {
	Name    string `json:"name"`
	Content string `json:"content"`
	Version int64  `json:"-"`
}

type TemplateModel struct {
	DB *sql.DB
}

func (m *TemplateModel) Insert(t *Template) error {
	query := `
INSERT INTO templates (name, content)
VAULES (?, ?)
RETURNING version`

	return m.DB.QueryRow(query, t.Name, t.Content).Scan(&t.Version)
}

func (m *TemplateModel) Get(name string) (*Template, error) {
	query := `
SELECT name, content, version
FROM templates
WHERE name = ?`

	var t Template

	err := m.DB.QueryRow(query, name).Scan(
		&t.Name,
		&t.Content,
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

func (m TemplateModel) Update(t *Template) error {
	// Template name should not be modified.
	query := `
UPDATE templates
SET content = ?, version = version + 1
WHERE name = ? AND version = ?
RETURNING version`

	err := m.DB.QueryRow(query, t.Content, t.Name, t.Version).Scan(&t.Version)
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

func (m TemplateModel) Delete(name string) error {
	query := `
DELETE FROM templates
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
