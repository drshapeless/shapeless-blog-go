package rest

import (
	"net/http"
	"time"
)

func (app *Application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Secret string `json:"secret"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Secret != app.Secret {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.Models.Tokens.New(24 * time.Hour)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSONInterface(w, http.StatusCreated, token, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
