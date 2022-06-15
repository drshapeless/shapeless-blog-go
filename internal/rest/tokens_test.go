package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

func TestToken(t *testing.T) {
	app := newApplication()

	// Create token
	tokenInput := struct {
		Secret string `json:"secret"`
	}{
		Secret: app.Secret,
	}

	tj, err := json.Marshal(&tokenInput)
	if err != nil {
		t.Fatal(err)
	}

	reader := bytes.NewReader(tj)

	req, err := http.NewRequest("POST", "/api/tokens/authentication", reader)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.createAuthenticationTokenHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf(
			"handler returned wrong status code: got %v want %v",
			status,
			http.StatusCreated,
		)
	}

	var token data.Token
	err = json.Unmarshal(rr.Body.Bytes(), &token)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(token)
}
