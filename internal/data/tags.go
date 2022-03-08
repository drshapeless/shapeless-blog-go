package data

import (
	"database/sql"
)

type Tag struct {
	PostID int64  `json:"post_id"`
	Tag    string `json:"tag"`
}

type TagModel struct {
	DB *sql.DB
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
		var tag string
		err := rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func (m *TagModel) GetPostsWithTag(tag string) ([]*Post, error) {
	query := `
SELECT posts.id, posts.title, posts.filename, posts.created, posts.updated
FROM posts
INNER JOIN tags
ON posts.id = tags.post_id
WHERE tags.tag = ?
ORDER BY posts.id DESC`

	rows, err := m.DB.Query(query, tag)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		var post Post
		err = rows.Scan(
			&post.ID,
			&post.Title,
			&post.Filename,
			&post.Created,
			&post.Updated,
		)

		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, err
}

func (m *TagModel) GetTagsWithPost(post *Post) ([]string, error) {
	query := `
SELECT tag
FROM tags
WHERE post_id = ?
`

	rows, err := m.DB.Query(query, post.ID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tags []string

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, err
}
