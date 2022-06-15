package rest

import (
	"errors"
	"net/http"

	"github.com/drshapeless/shapeless-blog/internal/data"
	"github.com/go-chi/chi"
)

func (app *Application) showTemplateHandler(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "title")

	temp, err := app.Models.Templates.GetByName(t)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSONInterface(w, http.StatusOK, temp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) createTemplateHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	temp := &data.Template{
		Name:    input.Name,
		Content: input.Content,
	}

	err = app.Models.Templates.Insert(temp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.updateTemplateCache(temp)

	err = app.writeJSONInterface(w, http.StatusCreated, temp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) updateTemplateHandler(w http.ResponseWriter, r *http.Request) {
	ti := chi.URLParam(r, "title")

	var input struct {
		Content string `json:"content"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	t, err := app.Models.Templates.GetByName(ti)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	t.Content = input.Content

	err = app.Models.Templates.Update(t)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	app.updateTemplateCache(t)

	err = app.writeJSONInterface(w, http.StatusOK, t, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) deleteTemplateHandler(w http.ResponseWriter, r *http.Request) {
	ti := chi.URLParam(r, "title")

	err := app.Models.Templates.Delete(ti)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
