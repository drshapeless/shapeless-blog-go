package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

func TestPosts(t *testing.T) {
	app := ExampleApplication()

	// Creating post
	postInput := struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		URL      string `json:"url"`
		CreateAt string `json:"create_at"`
		UpdateAt string `json:"update_at"`
	}{
		Title:    "Unit Testing",
		Content:  "Here is some unit testing content.",
		URL:      "unit-testing",
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
	handler := http.HandlerFunc(app.createPostHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusCreated,
		)
	}

	var post data.Post
	err = json.Unmarshal(rr.Body.Bytes(), &post)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(post)

	// Update post
	postInput.Content = "Here is some updated content."

	pj, err = json.Marshal(&postInput)
	if err != nil {
		t.Fatal(err)
	}

	reader = bytes.NewReader(pj)

	req, err = http.NewRequest("PATCH", "/api/bloggin/posts/id/1", reader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	q := req.URL.Query()
	q.Add("id", "1")
	req.URL.RawQuery = q.Encode()

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(app.updatePostHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Log(rr.Body.String())
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusOK,
		)
	}

	err = json.Unmarshal(rr.Body.Bytes(), &post)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(post)

}
