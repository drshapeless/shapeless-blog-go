package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/authentication/", app.createAuthenticationTokenHandler)

		r.Route("/posts", func(r chi.Router) {
			r.Use(app.authenticate)
			r.Post("/", app.createPostHandler)
			r.Get("/:id", app.showPostHandler)
			// Using patch here is not that accurate.
			// But fuck that.
			r.Patch("/:id", app.updatePostHandler)
			r.Delete("/:id", app.deletePostHandler)
		})

		r.Route("/templates", func(r chi.Router) {
			r.Use(app.authenticate)
			r.Post("/", app.createTemplateHandler)
			r.Get("/:name", app.showTemplateHandler)
			r.Patch("/:name", app.updateTemplateHandler)
			r.Delete("/:name", app.deleteTemplateHandler)
		})
	})

	r.Get("/", app.showHomeHTMLHandler)
	r.Get("/p/:id", app.showPostHTMLHandler)

	return r
}
