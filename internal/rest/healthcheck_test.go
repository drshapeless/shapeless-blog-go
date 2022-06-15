package rest

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
)

func TestHealthcheck(t *testing.T) {
	app := newApplication()

	t.Log("Testing healthcheck...")

	r := chi.NewRouter()
	r.Get("/api/healthcheck", app.healthcheckHandler)

	req, err := http.NewRequest("GET", "/api/healthcheck", nil)
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

	var o struct {
		Status  string `json:"status"`
		Version string `json:"version"`
	}
	t.Log(rr.Body.String())
	err = json.Unmarshal(rr.Body.Bytes(), &o)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(o)
}
