package rest

import (
	"net/http"

	"github.com/go-chi/chi"
)

// @title        shapeless-blog API
// @version      4.0.0
// @description  shapeless-blog api server.

// @host      https://blog.drshapeless.com
// @BasePath  /api
func (app *Application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", app.showHomeWebHandler)
	r.Get("/p/{title}.html", app.showPostWebHandler)
	r.Get("/t/{tag}.html", app.showTagWebHandler)

	r.Route("/api", func(r chi.Router) {
		r.Use(enableCORS)
		r.Get("/healthcheck", app.healthcheckHandler)

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
	r.Put("/posts/id/{id}", app.updatePostHandler)
	r.Delete("/posts/id/{id}", app.deletePostHandler)

	r.Post("/templates", app.createTemplateHandler)
	r.Get("/templates/{title}", app.showTemplateHandler)
	r.Put("/templates/{title}", app.updateTemplateHandler)
	r.Delete("/templates/{title}", app.deleteTemplateHandler)

	r.Get("/tags/{tag}", app.showTagHandler)
	r.Delete("/tags/{tag}", app.deleteTagHandler)
}

func (app *Application) tokenRoutes(r chi.Router) {
	r.Post("/authentication", app.createAuthenticationTokenHandler)
}
