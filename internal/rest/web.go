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
	Preview  string        `json:"preview"`
	Tags     []string      `json:"tags"`
	Content  template.HTML `json:"content"`
	CreateAt string        `json:"create_at"`
	UpdateAt string        `json:"update_at"`
}

func makeHtmlPost(p *data.Post, tags []string) *htmlPost {
	o := &htmlPost{
		ID:       p.ID,
		Title:    p.Title,
		URL:      p.URL,
		Preview:  p.Preview,
		Tags:     tags,
		Content:  template.HTML(p.Content),
		CreateAt: p.CreateAt,
		UpdateAt: p.UpdateAt,
	}

	return o
}

func (app *Application) showHomeWebHandler(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Posts []*htmlPost
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

	for _, post := range posts {
		tags, err := app.Models.Tags.GetTagsWithPostID(post.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		hpost := makeHtmlPost(post, tags)
		body.Posts = append(body.Posts, hpost)
	}

	tmpl := app.TemplateCache["home"]
	if tmpl == nil {
		app.emptyTemplateResponse(w, r)
		return
	}

	tmpl.Execute(w, body)
}

func (app *Application) showPostWebHandler(w http.ResponseWriter, r *http.Request) {
	// No need to show preview when showing post.
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

	tmpl := app.TemplateCache["post"]
	if tmpl == nil {
		app.emptyTemplateResponse(w, r)
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

	tmpl := app.TemplateCache["tag"]
	if tmpl == nil {
		app.emptyTemplateResponse(w, r)
		return
	}

	tmpl.Execute(w, body)
}
