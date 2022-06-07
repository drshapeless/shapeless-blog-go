package data_test

import (
	"testing"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

func TestGetAll(t *testing.T) {
	db, _ := data.OpenDB("/Users/jacky/shapeless-blog/shapeless-blog.db")
	m := data.NewModels(db)
	rows, err := m.Posts.GetAll()
	if err != nil {
		t.Error(err)
	}

	for _, v := range rows {
		t.Logf(v.Title)
	}
}
