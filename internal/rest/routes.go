package rest

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *Application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", app.showHomeHandler)
	r.Get("/posts/{title}", app.showPostHandler)
	r.Get("/tags/{tag}", app.showTagHandler)

	r.Route("/api", func(r chi.Router) {
		r.Use(enableCORS)
		r.Route("/tokens", app.tokenRoutes)

		r.Route("/posts", app.postRoutes)
	})

	return r
}

func (app *Application) postRoutes(r chi.Router) {
	r.Use(app.authenticate)

	r.Post("/create-post", app.createPostHandler)
	r.Post("/create-template", app.createTemplateHandler)
}

func (app *Application) tokenRoutes(r chi.Router) {
	r.Post("/authentication", app.createAuthenticationTokenHandler)
}
