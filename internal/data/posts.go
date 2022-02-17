package data

import (
	"database/sql"
	"errors"

	"github.com/pelletier/go-toml/query"
)

// Since there is no array in SQLite, we store the csv value instead.
// Also, there is no time.
type Post struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Tags      string `json:"tags"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Version   int    `json:"-"`
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(p *Post) error {
	query := `
INSERT INTO posts (title, content, tags, created_at, updated_at)
VAULES (?, ?, ?, ?, ?)
RETURNING id, version`

	args := []interface{}{p.Title, p.Content, p.Tags, p.CreatedAt, p.UpdatedAt}

	return m.DB.QueryRow(query, args...).Scan(&p.ID, &p.Version)
}

func (m *PostModel) Get(id int64) (*Post, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
SELECT id, title, content, tags, created_at, updated_at, version
FROM posts
WHERE id = ?`

	var p Post

	err := m.DB.QueryRow(query, id).Scan(
		&p.ID,
		&p.Title,
		&p.Content,
		&p.Tags,
		&p.CreatedAt,
		&p.UpdatedAt,
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

func (m PostModel) GetLatest() ([]*Post, Metadata, error) {
	query := ``
}

func (m PostModel) GetAll(filters Filters) ([]*Post, Metadata, error) {
	query := `
SELECT count(*) OVER(), id, title, content, tags, created_at, updated_at, version
FROM posts
LIMIT ? OFFSET ?`

	rows, err := m.DB.Query(query, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	posts := []*Post{}

	for rows.Next() {
		var p Post
		err := rows.Scan(
			&totalRecords,
			&p.ID,
			&p.Title,
			&p.Content,
			&p.Tags,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.Version,
		)

		if err != nil {
			return nil, Metadata{}, err
		}
		posts = append(posts, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return posts, metadata, nil
}

func (m PostModel) GetTag(id int64) (string, error) {
	query := `
SELECT tags
FROM posts
WHERE id = ?`

	var tags string
	err := m.DB.QueryRow(query, id).Scan(&tags)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", err
		}
	}

	return tags, nil
}

// This function is for extreme performance concern.
func (m PostModel) UpdateNoTags(p *Post) error {
	query := `
UPDATE posts
SET title = ?, content = ?, created_at = ?, updated_at = ?, version = version + 1
WHERE id = ? AND version = ?
RETURNING version`

	args := []interface{}{
		p.Title,
		p.Content,
		p.CreatedAt,
		p.UpdatedAt,
		p.ID,
		p.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&p.Version)
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

func (m PostModel) Update(p *Post) error {
	query := `
UPDATE posts
SET title = ?, content = ?, tags = ?, created_at = ?, updated_at = ?, version = version + 1
WHERE id = ? AND version = ?
RETURNING version`

	args := []interface{}{
		p.Title,
		p.Content,
		p.Tags,
		p.CreatedAt,
		p.UpdatedAt,
		p.ID,
		p.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&p.Version)
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

func (m PostModel) Delete(id int64) error {
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
