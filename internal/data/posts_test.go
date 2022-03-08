package data_test

import (
	"database/sql"
	"testing"

	"github.com/drshapeless/shapeless-blog/internal/data"
	_ "github.com/mattn/go-sqlite3"
)

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestGetAll(t *testing.T) {
	db, _ := openDB("/Users/jacky/shapeless-blog/shapeless-blog.db")
	m := data.NewModels(db)
	rows, err := m.Posts.GetAll()
	if err != nil {
		t.Error(err)
	}

	for _, v := range rows {
		t.Logf(v.Title)
	}
}
