package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestShowHomeWeb(t *testing.T) {
	app := newApplication()

	t.Log("Testing show home web...")

	r := chi.NewRouter()
	r.Get("/", app.showHomeWebHandler)

	req, err := http.NewRequest("GET", "/", nil)
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

	t.Log(rr.Body.String())
}

func TestShowPostWeb(t *testing.T) {
	app := newApplication()

	t.Log("Testing show post web...")

	r := chi.NewRouter()
	r.Get("/p/{title}", app.showPostWebHandler)

	req, err := http.NewRequest("GET", "/p/test", nil)
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

	t.Log(rr.Body.String())
}

func TestShowTagWeb(t *testing.T) {
	app := newApplication()

	t.Log("Testing show tag web...")

	r := chi.NewRouter()
	r.Get("/t/{tag}", app.showTagWebHandler)

	req, err := http.NewRequest("GET", "/t/test", nil)
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

	t.Log(rr.Body.String())
}
