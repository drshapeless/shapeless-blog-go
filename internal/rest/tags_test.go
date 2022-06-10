package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drshapeless/shapeless-blog/internal/data"
	"github.com/go-chi/chi"
)

func TestShowTag(t *testing.T) {
	app := ExampleApplication()

	t.Log("Testing show tag...")

	r := chi.NewRouter()
	r.Route("/api/blogging", func(r chi.Router) {
		r.Get("/tags/{tag}", app.showTagHandler)
	})

	req, err := http.NewRequest("GET", "/api/blogging/tags/test", nil)
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

	type ps []data.Post
	var pss ps
	t.Log(rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &pss)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(pss)
}

func TestDeleteTag(t *testing.T) {
	app := ExampleApplication()

	t.Log("Testing show tag...")

	r := chi.NewRouter()
	r.Route("/api/blogging", func(r chi.Router) {
		r.Delete("/tags/{tag}", app.deleteTagHandler)
	})

	req, err := http.NewRequest("DELETE", "/api/blogging/tags/test", nil)
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
