package data

import (
	"database/sql"
	"time"
)

// Since there is no array in SQLite, we store the csv value instead.
type Post struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      string    `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Version   int       `json:"-"`
}

type PostModel struct {
	DB *sql.DB
}
