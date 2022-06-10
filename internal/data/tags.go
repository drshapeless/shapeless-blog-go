package data

import "database/sql"

type Tag struct {
	PostID int    `json:"post_id"`
	Tag    string `json:"tag"`
}

type TagModel struct {
	DB *sql.DB
}

func (m *TagModel) Insert(t *Tag) error {
	query := `
INSERT INTO tags (post_id, tag)
VALUES (?, ?)`

	args := []interface{}{
		t.PostID,
		t.Tag,
	}

	result, err := m.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return err
	}

	return nil
}

func (m *TagModel) GetAllDistinctTags() ([]string, error) {
	query := `
SELECT DISTINCT tag
FROM tags
ORDER BY tag ASC
`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []string
	for rows.Next() {
		var t string
		err := rows.Scan(
			&t,
		)
		if err != nil {
			return nil, err
		}
		tags = append(tags, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (m *TagModel) GetPostsWithTag(t string, pagesize, page int) ([]*Post, error) {
	if pagesize < 1 || page < 1 {
		return nil, ErrRecordNotFound
	}

	// This will not return post content.
	query := `
SELECT posts.id, posts.title, posts.url, posts.create_at, posts.update_at
FROM posts
INNER JOIN tags
ON posts.id = tags.post_id
WHERE tags.tag = ?
ORDER BY posts.id DESC
LIMIT ? OFFSET ?`

	args := []interface{}{
		t,
		pagesize,
		calculateOffset(pagesize, page),
	}

	rows, err := m.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		var p Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.URL,
			&p.CreateAt,
			&p.UpdateAt,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, err
}

func (m *TagModel) GetTagsWithPostID(pid int) ([]string, error) {
	query := `
SELECT tag
FROM tags
WHERE post_id = ?
`

	rows, err := m.DB.Query(query, pid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []string

	for rows.Next() {
		var t string
		err = rows.Scan(
			&t,
		)
		if err != nil {
			return nil, err
		}

		tags = append(tags, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, err
}

func (m TagModel) Delete(tag string, pid int) error {
	if pid < 1 {
		return ErrRecordNotFound
	}

	query := `
DELETE FROM tags
WHERE tag = ? AND post_id = ?`

	result, err := m.DB.Exec(query, tag, pid)
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

func (m TagModel) DeleteAllWithPostID(pid int) error {
	if pid < 1 {
		return ErrRecordNotFound
	}

	query := `
DELETE from tags
WHERE post_id = ?`

	result, err := m.DB.Exec(query, pid)
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

func (m TagModel) DeleteAllForTag(tag string) error {
	if tag == "" {
		return ErrRecordNotFound
	}

	query := `
DELETE FROM tags
WHERE tag = ?
`

	result, err := m.DB.Exec(query, tag)
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
