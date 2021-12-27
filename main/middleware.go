package main

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (app *application) requireAuthentication(next httprouter.Handle) httprouter.Handle {
	return httprouter.Handle(func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Add("Vary", "authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			app.authenticationRequiredResponse(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		err := app.compareToken(token)
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		next(w, r, ps)
	})
}
