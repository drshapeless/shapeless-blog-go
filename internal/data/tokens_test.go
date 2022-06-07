package data_test

import (
	"testing"
	"time"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

func TestToken(t *testing.T) {
	db, err := data.OpenDB("shapeless-blog.db")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Resetting test database.")
	data.ResetDB(db)
	t.Log("Database reset.")

	m := data.NewModels(db)

	t.Log("Generating new token")
	tok, err := m.Tokens.New(0)
	if err != nil {
		t.Error(err)
	}
	t.Log(tok)

	time.Sleep(time.Second)

	t.Log("Deleting expired tokens")
	err = m.Tokens.DeleteAllForExpired()
	if err != nil {
		t.Error(err)
	}
}
