package data

import (
	"database/sql"
	"errors"
)

// Filename is just the lower case of the title with spaces replaced
// by underline.
type Post struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Filename string `json:"filename"`
	Created  string `json:"created"`
	Updated  string `json:"updated"`
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) GetWithID(id int64) (*Post, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
SELECT id, title, filename, created, updated
FROM posts
WHERE id = ?`

	var p Post

	err := m.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Title,
		&p.Filename,
		&p.Created,
		&p.Updated,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &p, nil
}

func (m PostModel) GetWithFilename(name string) (*Post, error) {
	query := `
SELECT id, title, filename, created, updated
FROM posts
WHERE filename = ?`

	var p Post

	err := m.DB.QueryRow(query, name).Scan(
		&p.ID,
		&p.Title,
		&p.Filename,
		&p.Created,
		&p.Updated,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &p, nil
}

func (m PostModel) GetAll() ([]*Post, error) {
	query := `
SELECT id, title, filename, created, updated
FROM posts
ORDER BY id DESC`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	posts := []*Post{}

	for rows.Next() {
		var p Post
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Filename,
			&p.Created,
			&p.Updated,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
