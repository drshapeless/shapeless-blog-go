package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/drshapeless/shapeless-blog/internal/data"
	"github.com/drshapeless/shapeless-blog/internal/validator"
)

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func enableCORS(next http.Handler) http.Handler {
	// Enabling cors is a good to have function.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://drshapeless.com")

		if r.Method == http.MethodOptions &&
			r.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set(
				"Access-Control-Allow-Methods",
				"GET, POST, PUT, OPTIONS",
			)
			w.Header().Set(
				"Access-Control-Allow-Headers",
				"Authorization, X-PINGOTHER, Content-Type",
			)
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (app *Application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			app.authenticationRequiredResponse(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidCredentialsResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		_, err := app.Models.Tokens.GetByPlaintext(token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidCredentialsResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		next.ServeHTTP(w, r)
	})
}
