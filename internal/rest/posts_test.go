package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/drshapeless/shapeless-blog/internal/data"
	"github.com/go-chi/chi"
)

func TestCreatePost(t *testing.T) {
	app := newApplication()

	t.Log("Testing create post...")

	r := chi.NewRouter()
	r.Route("/api/blogging", func(r chi.Router) {
		r.Post("/posts", app.createPostHandler)
	})

	// Creating post
	postInput := struct {
		Title    string   `json:"title"`
		URL      string   `json:"url"`
		Tags     []string `json:"tags"`
		Content  string   `json:"content"`
		CreateAt string   `json:"create_at"`
		UpdateAt string   `json:"update_at"`
	}{
		Title:    "Unit Testing",
		URL:      "unit-testing",
		Tags:     []string{"unit", "test"},
		Content:  "Here is some unit testing content.",
		CreateAt: time.Now().Format(dateLayout),
		UpdateAt: time.Now().Add(time.Hour).Format(dateLayout),
	}

	pj, err := json.Marshal(&postInput)
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(pj)

	req, err := http.NewRequest("POST", "/api/blogging/posts", reader)
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

	var post restPost
	t.Log(rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &post)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(post)
}

func TestUpdatePost(t *testing.T) {
	app := newApplication()

	t.Log("Testing update post...")

	r := chi.NewRouter()
	r.Route("/api/blogging", func(r chi.Router) {
		r.Put("/posts/id/{id}", app.updatePostHandler)
	})

	pi := struct {
		Tags    []string `json:"tags"`
		Content string   `json:"content"`
	}{
		Tags:    []string{"update", "test"},
		Content: "Here is some updated content.",
	}

	pj, err := json.Marshal(&pi)
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(pj)

	req, err := http.NewRequest("PUT", "/api/blogging/posts/id/1", reader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK,
		)
	}

	var post restPost
	t.Log(rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &post)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(post)
}

func TestShowPostWithID(t *testing.T) {
	app := newApplication()

	t.Log("Testing show post with id...")

	r := chi.NewRouter()
	r.Route("/api/blogging", func(r chi.Router) {
		r.Get("/posts/id/{id}", app.showPostWithIDHandler)
	})

	req, err := http.NewRequest("GET", "/api/blogging/posts/id/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK,
		)
	}

	var post restPost
	t.Log(rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &post)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(post)
}

func TestShowPostWithTitle(t *testing.T) {
	app := newApplication()

	t.Log("Testing show post with title...")

	r := chi.NewRouter()
	r.Route("/api/blogging", func(r chi.Router) {
		r.Get("/posts/{title}", app.showPostWithTitleHandler)
	})

	req, err := http.NewRequest("GET", "/api/blogging/posts/unit-testing", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK,
		)
		return
	}

	var post restPost
	t.Log(rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &post)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(post)
}

func TestDeletePost(t *testing.T) {
	app := newApplication()

	t.Log("Testing delete post...")

	r := chi.NewRouter()
	r.Route("/api/blogging", func(r chi.Router) {
		r.Delete("/posts/id/{id}", app.deletePostHandler)
	})

	req, err := http.NewRequest("DELETE", "/api/blogging/posts/id/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Log(rr.Body.String())
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusNoContent,
		)
		return
	}
}

func TestPostAll(t *testing.T) {
	p := os.Getenv("SHAPELESS_BLOG_DB_PATH")
	db, err := data.OpenDB(p)
	if err != nil {
		t.Fatal(err)
	}

	data.ResetDB(db)

	TestCreatePost(t)
	TestUpdatePost(t)
	TestShowPostWithID(t)
	TestShowPostWithTitle(t)
	TestDeletePost(t)
}
