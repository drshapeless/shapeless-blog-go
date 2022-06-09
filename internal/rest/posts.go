package rest

import (
	"errors"
	"net/http"

	"github.com/drshapeless/shapeless-blog/internal/data"
	"github.com/go-chi/chi"
)

func (app *Application) showPostWithTitleHandler(w http.ResponseWriter, r *http.Request) {
	tt := chi.URLParam(r, "title")
	p, err := app.Models.Posts.GetWithURL(tt)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSONInterface(w, http.StatusOK, p, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) showPostWithIDHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	p, err := app.Models.Posts.GetWithID(int(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSONInterface(w, http.StatusOK, p, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		URL      string `json:"url"`
		CreateAt string `json:"create_at"`
		UpdateAt string `json:"update_at"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	p := &data.Post{
		Title:    input.Title,
		Content:  input.Content,
		URL:      input.URL,
		CreateAt: input.CreateAt,
		UpdateAt: input.UpdateAt,
	}

	err = app.Models.Posts.Insert(p)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSONInterface(w, http.StatusCreated, p, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Title    *string `json:"title"`
		Content  *string `json:"content"`
		URL      *string `json:"url"`
		CreateAt *string `json:"create_at"`
		UpdateAt *string `json:"update_at"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	p, err := app.Models.Posts.GetWithID(int(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if input.Title != nil {
		p.Title = *input.Title
	}

	if input.Content != nil {
		p.Content = *input.Content
	}

	if input.URL != nil {
		p.URL = *input.URL
	}

	if input.CreateAt != nil {
		p.CreateAt = *input.CreateAt
	}

	if input.UpdateAt != nil {
		p.UpdateAt = *input.UpdateAt
	}

	err = app.Models.Posts.Update(p)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSONInterface(w, http.StatusOK, p, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.Models.Posts.Delete(int(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
