package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drshapeless/shapeless-blog/internal/data"
	"github.com/go-chi/chi"
)

func TestCreateTemplate(t *testing.T) {
	app := ExampleApplication()

	t.Log("Testing create template...")

	r := chi.NewRouter()
	r.Post("/api/blogging/templates", app.createTemplateHandler)

	input := struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}{
		Name:    "test",
		Content: "test template content",
	}

	i, err := json.Marshal(&input)
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(i)

	req, err := http.NewRequest("POST", "/api/blogging/templates", reader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := r

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusCreated,
		)
	}

	var tem data.Template
	t.Log(rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &tem)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tem)
}

func TestShowTemplate(t *testing.T) {
	app := ExampleApplication()

	t.Log("Testing show template...")

	r := chi.NewRouter()
	r.Get("/api/blogging/templates/{title}", app.showTemplateHandler)

	req, err := http.NewRequest("GET", "/api/blogging/templates/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := r

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK,
		)
	}

	var tem data.Template
	t.Log(rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &tem)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tem)
}

func TestUpdateTemplate(t *testing.T) {
	app := ExampleApplication()

	t.Log("Testing update template...")

	r := chi.NewRouter()
	r.Patch("/api/blogging/templates/{title}", app.updateTemplateHandler)

	input := struct {
		Content string `json:"content"`
	}{
		Content: "updated template content",
	}

	i, err := json.Marshal(&input)
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(i)

	req, err := http.NewRequest("PATCH", "/api/blogging/templates/test", reader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := r

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK,
		)
	}

	var tem data.Template
	t.Log(rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &tem)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tem)
}

func TestDeleteTemplate(t *testing.T) {
	app := ExampleApplication()

	t.Log("Testing delete template...")

	r := chi.NewRouter()
	r.Delete("/api/blogging/templates/{title}", app.deleteTemplateHandler)

	req, err := http.NewRequest("DELETE", "/api/blogging/templates/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := r

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusNoContent,
		)
	}
}
