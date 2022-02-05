package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/authentication/", app.createAuthenticationTokenHandler)

		r.Get("/posts/:id", app.showPostHandler)
		r.Patch("/posts/:id", app.updatePostHandler)
		r.Delete("/posts/:id", app.deletePostHandler)
		r.Post("/posts/", app.createPostHandler)
	})

	return r
}
