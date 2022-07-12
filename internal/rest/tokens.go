package rest

import (
	"net/http"
	"time"
)

type secretInput struct {
	Secret string `json:"secret"`
}

// createAuthenticationTokenHandler
// @Summary  Create authentication token
// @Description
// @Tags     tokens
// @Accept   json
// @Produce  json
// @Param    data  body      secretInput  true  "Secret object"
// @Success  201   {object}  data.Token
// @Failure  400   {object}  errorObject
// @Failure  401   {object}  errorObject
// @Failure  500   {object}  errorObject
// @Router   /blogging/tokens [post]
func (app *Application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input secretInput

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Secret != app.Secret {
		app.invalidCredentialsResponse(w, r)
		return
	}

	app.Models.Tokens.DeleteAllForExpired()

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
