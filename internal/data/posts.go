package data

import (
	"database/sql"
	"errors"
	"fmt"
)

type Post struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Content  string `json:"content"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
	Version  int    `json:"-"`
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(p *Post) error {
	query := `
INSERT INTO posts (title, content, url, create_at, update_at)
VALUES (?, ?, ?, ?, ?)
RETURNING id, version
`
	args := []interface{}{
		p.Title,
		p.Content,
		p.URL,
		p.CreateAt,
		p.UpdateAt,
	}

	err := m.DB.QueryRow(query, args...).Scan(
		&p.ID,
		&p.Version,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *PostModel) getWith(keyword string, value interface{}, ty int) (*Post, error) {
	query := `
SELECT id, title, url, content, create_at, update_at, version
FROM posts
`

	query += fmt.Sprintf("WHERE %s = ?", keyword)

	var p Post

	var tmp *sql.Row
	switch ty {
	case 0:
		tmp = m.DB.QueryRow(query, value.(int))
	case 1:
		tmp = m.DB.QueryRow(query, value.(string))
	default:
		return nil, ErrRecordNotFound
	}

	err := tmp.Scan(
		&p.ID,
		&p.Title,
		&p.URL,
		&p.Content,
		&p.CreateAt,
		&p.UpdateAt,
		&p.Version,
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

func (m *PostModel) GetWithID(id int) (*Post, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	return m.getWith("id", id, 0)
}

func (m PostModel) GetWithTitle(title string) (*Post, error) {
	return m.getWith("title", title, 1)
}

func (m PostModel) GetWithURL(url string) (*Post, error) {
	return m.getWith("url", url, 1)
}

func (m PostModel) GetAll(pagesize, page int) ([]*Post, error) {
	// This function will not return content.
	if pagesize < 1 || page < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
SELECT id, title, url, create_at, update_at, version
FROM posts
ORDER BY id DESC
LIMIT ? OFFSET ?
`
	args := []interface{}{
		pagesize,
		calculateOffset(pagesize, page),
	}

	rows, err := m.DB.Query(query, args...)
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
			&p.URL,
			&p.CreateAt,
			&p.UpdateAt,
			&p.Version,
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

func (m PostModel) Update(p *Post) error {
	query := `
UPDATE posts
SET title = ?, url = ?, content = ?, create_at = ?, update_at = ?, version = version + 1
WHERE id = ? AND version = ?
RETURNING version`

	args := []interface{}{
		p.Title,
		p.URL,
		p.Content,
		p.CreateAt,
		p.UpdateAt,
		p.ID,
		p.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(
		&p.Version,
	)

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

func (m PostModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
DELETE FROM posts
WHERE id = ?`

	result, err := m.DB.Exec(query, id)
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
