package main

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		Tags      []string `json:"tags"`
		CreatedAt string   `json:"created_at"`
		UpdatedAt string   `json:"updated_at"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Sort the slices first.
	sort.Strings(input.Tags)

	post := &data.Post{
		Title:     input.Title,
		Content:   input.Content,
		Tags:      slice2csv(input.Tags),
		CreatedAt: input.CreatedAt,
		UpdatedAt: input.UpdatedAt,
	}

	err = app.models.Posts.Insert(post)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Fuck with the tags here.
	for _, v := range input.Tags {
		if app.models.Tags.IsExist(v) {
			err = app.models.Tags.AddPostID(post.ID, v)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
		} else {
			t := &data.Tag{
				Name:   v,
				PostID: fmt.Sprintf("%d", post.ID),
			}
			err = app.models.Tags.Insert(t)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
		}
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/posts/%d", post.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"post": post}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	post, err := app.models.Posts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	post, err := app.models.Posts.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Every update request must contain all the information.
	// That's how the Emacs client works.
	var input struct {
		Title     string   `json:"title"`
		Content   string   `json:"content"`
		Tags      []string `json:"tags"`
		CreatedAt string   `json:"created_at"`
		UpdatedAt string   `json:"updated_at"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Check if the tags is modified.
	// If so, that is a fucking nightmare.
	sort.Strings(input.Tags)
	if slice2csv(input.Tags) != post.Tags {

	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {

}
