package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/drshapeless/shapeless-blog/internal/data"

	"github.com/go-chi/chi"
	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	port int
	dir  string
}

type application struct {
	config config
	models data.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.dir, "dir", os.Getenv("HOME")+"/shapeless-blog", "shapeless-blog directory")
	flag.IntVar(&cfg.port, "port", 9398, "shapeless-blog port")
	migrate := flag.Bool("migrate", false, "Migrate shapeless-blog database")

	flag.Parse()

	if *migrate {
		db, err := openDB(cfg.dir + "/shapeless-blog.db")
		if err != nil {
			log.Fatal(err)
		}
		data.Migrate(db)
		log.Printf("Successfully installed shapeless-blog database\n")
		return
	}

	var app application
	app.config = cfg

	if !isExist(app.config.dir) {
		fmt.Fprintf(os.Stderr, "shapeless-blog directory does not exist!\n")
		return
	}

	if !isExist(app.databasePath()) {
		fmt.Fprintf(os.Stderr, "database not found!\n")
		return
	} else {
		db, err := openDB(app.config.dir + "/shapeless-blog.db")
		if err != nil {
			log.Fatal(err)
		}
		app.models = data.NewModels(db)
	}

	app.serve()
}

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func isExist(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting shapeless-blog server at port %d\n", app.config.port)
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}
	log.Printf("stopped shapeless-blog server\n")

	return nil
}

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Get("/", app.showHomeHandler)
	r.Get("/p/{title}.html", app.showPostHandler)
	r.Get("/t/{tag}", app.showTagHandler)
	return r
}

func (app *application) showHomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(app.templatePath("home"))
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	var body struct {
		Posts []*data.Post
		Tags  []string
	}

	posts, err := app.models.Posts.GetAll()
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	body.Posts = posts

	tags, err := app.models.Tags.GetAllDistinctTags()
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	body.Tags = tags

	tmpl.Execute(w, body)
}

func (app *application) showPostHandler(w http.ResponseWriter, r *http.Request) {
	filename := chi.URLParam(r, "title")

	post, err := app.models.Posts.GetWithFilename(filename)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	content, err := os.ReadFile(app.postPath(filename))
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	var body struct {
		Post    *data.Post
		Content template.HTML
	}
	body.Post = post
	body.Content = template.HTML(content)

	tmpl, err := template.ParseFiles(app.templatePath("post"))
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	tmpl.Execute(w, body)
}

func (app *application) showTagHandler(w http.ResponseWriter, r *http.Request) {
	tag := chi.URLParam(r, "tag")

	posts, err := app.models.Tags.GetPostsWithTag(tag)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	var body struct {
		Posts []*data.Post
		Tag   string
	}
	body.Posts = posts
	body.Tag = tag

	tmpl, err := template.ParseFiles(app.templatePath("tag"))
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	tmpl.Execute(w, body)
}
