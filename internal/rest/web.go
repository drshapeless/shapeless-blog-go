package rest

import (
	"errors"
	"html/template"
	"net/http"

	"github.com/drshapeless/shapeless-blog/internal/data"
	"github.com/drshapeless/shapeless-blog/internal/validator"
	"github.com/go-chi/chi"
)

type htmlPost struct {
	ID       int           `json:"id"`
	Title    string        `json:"title"`
	URL      string        `json:"url"`
	Tags     []string      `json:"tags"`
	Content  template.HTML `json:"content"`
	CreateAt string        `json:"create_at"`
	UpdateAt string        `json:"update_at"`
}

func makeHtmlPost(p *data.Post, t []string) *htmlPost {
	o := &htmlPost{
		ID:       p.ID,
		Title:    p.Title,
		URL:      p.URL,
		Tags:     t,
		Content:  template.HTML(p.Content),
		CreateAt: p.CreateAt,
		UpdateAt: p.UpdateAt,
	}

	return o
}

func (app *Application) showHomeWebHandler(w http.ResponseWriter, r *http.Request) {
	ts, err := app.Models.Templates.GetByName("home")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tmpl, err := template.New("home").Parse(ts.Content)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	var body struct {
		Posts []*data.Post
		Tags  []string
	}

	v := validator.New()

	qs := r.URL.Query()

	page := app.readInt(qs, "page", 1, v)
	pagesize := app.readInt(qs, "pagesize", 100, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	posts, err := app.Models.Posts.GetAll(pagesize, page)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	body.Posts = posts

	tags, err := app.Models.Tags.GetAllDistinctTags()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	body.Tags = tags

	tmpl.Execute(w, body)
}

func (app *Application) showPostWebHandler(w http.ResponseWriter, r *http.Request) {
	ti := chi.URLParam(r, "title")
	post, err := app.Models.Posts.GetWithURL(ti)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	tags, err := app.Models.Tags.GetTagsWithPostID(post.ID)
	if err != nil && !errors.Is(err, data.ErrRecordNotFound) {
		app.serverErrorResponse(w, r, err)
		return
	}

	body := makeHtmlPost(post, tags)

	ts, err := app.Models.Templates.GetByName("post")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	tmpl, err := template.New("post").Parse(ts.Content)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tmpl.Execute(w, body)
}

func (app *Application) showTagWebHandler(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")

	v := validator.New()
	qs := r.URL.Query()

	page := app.readInt(qs, "page", 1, v)
	pagesize := app.readInt(qs, "pagesize", 100, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	ps, err := app.Models.Tags.GetPostsWithTag(tag, pagesize, page)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var body struct {
		Posts []*data.Post
		Tag   string
	}
	body.Posts = ps
	body.Tag = tag

	ts, err := app.Models.Templates.GetByName("tag")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tmpl, err := template.New("tag").Parse(ts.Content)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tmpl.Execute(w, body)
}
