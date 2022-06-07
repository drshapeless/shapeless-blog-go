package data

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Models struct {
	Posts     PostModel
	Tags      TagModel
	Tokens    TokenModel
	Templates TemplateModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Posts:     PostModel{DB: db},
		Tags:      TagModel{DB: db},
		Tokens:    TokenModel{DB: db},
		Templates: TemplateModel{DB: db},
	}
}

func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Migrate(db *sql.DB) error {
	query := `
CREATE TABLE IF NOT EXISTS posts(
id         INTEGER PRIMARY KEY,
title      TEXT NOT NULL,
url        TEXT NOT NULL,
content    TEXT NOT NULL,
create_at  TEXT NOT NULL,
update_at  TEXT NOT NULL,
version    INTEGER DEFAULT 1
);

CREATE TABLE IF NOT EXISTS tags(
post_id  INTEGER NOT NULL,
tag      TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens(
hash  TEXT NOT NULL,
expiry INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS templates(
name    TEXT NOT NULL,
content TEXT NOT NULL,
version INTEGER DEFAULT 1
);
`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func Drop(db *sql.DB) error {
	query := `
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS templates;
`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func ResetDB(db *sql.DB) error {
	err := Drop(db)
	if err != nil {
		return err
	}

	err = Migrate(db)
	if err != nil {
		return err
	}

	return nil
}
