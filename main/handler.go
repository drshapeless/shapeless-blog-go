package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/drshapeless/shapeless-blog/data"
	"github.com/julienschmidt/httprouter"
)

func (app *application) createBlogpostHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var input struct {
		Title    string   `json:"title"`
		Created  string   `json:"created"`
		Updated  string   `json:"updated"`
		Category []string `json:"category"`
		Content  string   `json:"content"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	blog := data.Blog{
		Content: input.Content,
		Metadata: data.PostMetadata{
			Created:  input.Created,
			Updated:  input.Updated,
			Category: input.Category,
			Title:    input.Title,
		},
	}

	if blog.Content == "" {
		app.failValidationResponse(w, r, "content cannot be empty")
		return
	}
	if blog.Metadata.Title == "" {
		app.failValidationResponse(w, r, "title cannot be empty")
		return
	}

	id, err := app.createPost(blog)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/%d", id))

	err = app.writeJSON(w, http.StatusCreated, envelope{"id": id}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showBlogpostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "the requested resource could not be found\n")
		return
	}

	blog, err := app.getBlog(id)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	tmpl := app.showBlogpostTemplate()
	tmpl.Execute(w, blog)
}

func (app *application) updateBlogpostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		ID       int      `json:"id"`
		Title    string   `json:"title"`
		Created  string   `json:"created"`
		Updated  string   `json:"updated"`
		Category []string `json:"category"`
		Content  string   `json:"content"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.ID == 0 {
		app.failValidationResponse(w, r, "must provided a valid id")
		return
	}
	if input.Title == "" {
		app.failValidationResponse(w, r, "title cannot be empty")
		return
	}
	if input.Content == "" {
		app.failValidationResponse(w, r, "content cannot be empty")
		return
	}

	blog := data.Blog{
		Content: input.Content,
		Metadata: data.PostMetadata{
			ID:       input.ID,
			Title:    input.Title,
			Created:  input.Created,
			Updated:  input.Updated,
			Category: input.Category,
		},
	}

	err = app.updatePost(blog)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"id": input.ID}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteBlogpostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "the requested resource could not be found\n")
		return
	}

	err = app.deletePost(id)
	if err != nil {
		fmt.Fprintf(w, "%s\n", err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "blog post %d successfully deleted"}, nil)
	// fmt.Fprintf(w, "blog post %d successfully deleted\n", id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getBlogpostHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	blog, err := app.getBlog(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"blogpost": blog}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showHomePageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	postids := app.cache.PostIDs

	sort.Sort(sort.Reverse(sort.IntSlice(postids)))

	var datas []data.PostMetadata
	for _, value := range postids {
		datas = append(datas, app.getMetadata(value))
	}

	tmpl := app.showHomePageTemplate()
	tmpl.Execute(w, datas)
}

func (app *application) showCategoriesIndexHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	tmpl := app.showCategoriesIndexTemplate()
	tmpl.Execute(w, app.cache.Categories)
}

func (app *application) showCategoryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("category")
	if !isStringExist(name, app.cache.Categories) {
		app.notFoundResponse(w, r)
		return
	}

	postids, err := readFileAsIntArray(app.categoryPath(name))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(postids)))

	var datas []data.PostMetadata

	for _, value := range postids {
		datas = append(datas, app.getMetadata(value))
	}

	var templateData struct {
		Category string
		Metadata []data.PostMetadata
	}

	templateData.Category = name
	templateData.Metadata = datas

	tmpl := app.showCategoryTemplate()
	tmpl.Execute(w, templateData)
}

func (app *application) authenticationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var input struct {
		User     string `json:"username"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.authenticateUser(input.User, input.Password)
	if err != nil {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token := generateToken()
	err = app.writeTokenHash(token)
	if err != nil {
		panic(err)
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"token": token}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
