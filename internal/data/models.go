package data

import "database/sql"

type Models struct {
	Posts PostModel
	Tags  TagModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Posts: PostModel{DB: db},
		Tags:  TagModel{DB: db},
	}
}

func Migrate(db *sql.DB) error {
	query := `
CREATE TABLE IF NOT EXISTS posts(
id INTEGER PRIMARY KEY,
title TEXT NOT NULL,
content TEXT NOT NULL,
tags TEXT NOT NULL,
created_at TEXT NOT NULL,
updated_at TEXT NOT NULL,
version INTEGER DEFAULT 1
);
CREATE TABLE IF NOT EXISTS tags(
name TEXT NOT NULL UNIQUE,
post_id TEXT NOT NULL
);
`

	// The result does not matter.
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
