package rest

import (
	"errors"
	"net/http"

	"github.com/drshapeless/shapeless-blog/internal/data"
	"github.com/go-chi/chi"
)

type restPost struct {
	ID       int      `json:"id"`
	Title    string   `json:"title"`
	URL      string   `json:"url"`
	Tags     []string `json:"tags"`
	Content  string   `json:"content"`
	CreateAt string   `json:"create_at"`
	UpdateAt string   `json:"update_at"`
}

func makeOutputPost(post *data.Post, tags []string) *restPost {
	o := &restPost{
		ID:       post.ID,
		Title:    post.Title,
		URL:      post.URL,
		Tags:     tags,
		Content:  post.Content,
		CreateAt: post.CreateAt,
		UpdateAt: post.UpdateAt,
	}

	return o
}

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

	tags, err := app.Models.Tags.GetTagsWithPostID(p.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	rp := makeOutputPost(p, tags)

	err = app.writeJSONInterface(w, http.StatusOK, rp, nil)
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

	tags, err := app.Models.Tags.GetTagsWithPostID(p.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	rp := makeOutputPost(p, tags)

	err = app.writeJSONInterface(w, http.StatusOK, rp, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input restPost

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

	for _, v := range input.Tags {
		t := data.Tag{
			PostID: p.ID,
			Tag:    v,
		}
		err = app.Models.Tags.Insert(&t)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	rp := makeOutputPost(p, input.Tags)

	err = app.writeJSONInterface(w, http.StatusCreated, rp, nil)
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
		Title    *string  `json:"title"`
		URL      *string  `json:"url"`
		Content  *string  `json:"content"`
		Tags     []string `json:"tags"`
		CreateAt *string  `json:"create_at"`
		UpdateAt *string  `json:"update_at"`
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

	if len(input.Tags) > 0 {
		err = app.Models.Tags.DeleteAllWithPostID(p.ID)
		if err != nil && !errors.Is(err, data.ErrRecordNotFound) {
			app.serverErrorResponse(w, r, err)
			return
		}
		for _, v := range input.Tags {
			t := data.Tag{
				PostID: p.ID,
				Tag:    v,
			}
			err = app.Models.Tags.Insert(&t)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
		}
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

	err = app.Models.Tags.DeleteAllWithPostID(int(id))
	if err != nil && !errors.Is(err, data.ErrRecordNotFound) {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
