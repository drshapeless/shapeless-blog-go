package data_test

import (
	"testing"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

func TestTemplate(t *testing.T) {
	db, err := data.OpenDB("shapeless-blog.db")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Resetting test database.")
	data.ResetDB(db)
	t.Log("Database reset.")

	m := data.NewModels(db)

	tem := data.Template{
		Name:    "test",
		Content: "test template",
	}

	t.Log("Inserting test template")
	err = m.Templates.Insert(&tem)
	if err != nil {
		t.Error(err)
	}

	t.Log("Updating test template")
	tem.Content = "updated test template"
	err = m.Templates.Update(&tem)
	if err != nil {
		t.Error(err)
	}

	t.Log("Getting test template")
	temp, err := m.Templates.GetByName("test")
	if err != nil {
		t.Error(err)
	}
	t.Log(temp)

	t.Log("Deleting test template")
	err = m.Templates.Delete("test")
	if err != nil {
		t.Error(err)
	}
}
