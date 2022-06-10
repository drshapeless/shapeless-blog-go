package rest

import (
	"net/http"

	"github.com/drshapeless/shapeless-blog/internal/validator"
	"github.com/go-chi/chi"
)

func (app *Application) showTagHandler(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")
	qs := r.URL.Query()

	v := validator.New()

	page := app.readInt(qs, "page", 1, v)
	pagesize := app.readInt(qs, "page_size", 20, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	ps, err := app.Models.Tags.GetPostsWithTag(tag, pagesize, page)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSONInterface(w, http.StatusOK, ps, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) deleteTagHandler(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")

	err := app.Models.Tags.DeleteAllForTag(tag)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
