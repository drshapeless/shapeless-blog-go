package main

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

func (app *application) showHomeHTMLHandler(w http.ResponseWriter, r *http.Request) {
	posts, err := app.models.Posts.GetAll()

	t, err := app.models.Templates.Get("home")
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	tmpl, err := template.ParseGlob(t.Content)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tmpl.Execute(w, posts)
}
