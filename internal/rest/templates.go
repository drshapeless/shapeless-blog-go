package rest

import (
	"errors"
	"net/http"

	"github.com/drshapeless/shapeless-blog/internal/data"
	"github.com/go-chi/chi"
)

// showTemplateHandler
// @Summary  Show template
// @Description
// @Tags     templates
// @Produce  json
// @Param    Authorization  header    string  true  "Bearer"
// @Param    title          body      string  true  "Template title string"
// @Success  200            {object}  data.Template
// @Failure  404            {object}  errorObject
// @Failure  500            {object}  errorObject
// @Router   /blogging/templates/{title} [get]
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

// createTemplateHandler
// @Summary  Create template
// @Description
// @Tags     templates
// @Accept   json
// @Produce  json
// @Param    Authorization  header    string         true  "Bearer"
// @Param    data           body      data.Template  true  "Template object"
// @Success  201            {object}  data.Template
// @Failure  400            {object}  errorObject
// @Failure  500            {object}  errorObject
// @Router   /blogging/templates [post]
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

// updateTemplateHandler
// @Summary  Update template
// @Description
// @Tags     templates
// @Accept   json
// @Produce  json
// @Param    Authorization  header    string         true  "Bearer"
// @Param    data           body      data.Template  true  "Template object"
// @Param    title          path      string         true  "Template title"
// @Success  200            {object}  data.Template
// @Failure  400            {object}  errorObject
// @Failure  404            {object}  errorObject
// @Failure  409            {object}  errorObject
// @Failure  500            {object}  errorObject
// @Router   /blogging/templates/{title} [put]
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

// deleteTemplateHandler
// @Summary  Delete template
// @Description
// @Tags     templates
// @Accept   json
// @Produce  json
// @Param    Authorization  header  string  true  "Bearer"
// @Param    title          body    string  true  "Template title"
// @Success  204            "No content"
// @Failure  409            {object}  errorObject
// @Failure  500            {object}  errorObject
// @Router   /blogging/templates/{title} [delete]
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
