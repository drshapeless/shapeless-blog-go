package rest

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *Application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", app.showHomeWebHandler)
	r.Get("/posts/{title}", app.showPostWebHandler)
	r.Get("/tags/{tag}", app.showTagWebHandler)

	r.Route("/api", func(r chi.Router) {
		r.Use(enableCORS)
		r.Route("/tokens", app.tokenRoutes)

		r.Route("/blogging", app.bloggingRoutes)
	})

	return r
}

func (app *Application) bloggingRoutes(r chi.Router) {
	r.Use(app.authenticate)

	r.Post("/posts", app.createPostHandler)
	r.Get("/posts/{title}", app.showPostWithTitleHandler)
	r.Get("/posts/id/{id}", app.showPostWithIDHandler)
	r.Patch("/posts/id/{id}", app.updatePostHandler)
	r.Delete("/posts/id/{id}", app.deletePostHandler)

	r.Post("/templates", app.createTemplateHandler)
	r.Get("/templates/{title}", app.showTemplateHandler)
	r.Patch("/templates/{title}", app.updateTemplateHandler)
	r.Delete("/templates/{title}", app.deleteTemplateHandler)
}

func (app *Application) tokenRoutes(r chi.Router) {
	r.Post("/authentication", app.createAuthenticationTokenHandler)
}
