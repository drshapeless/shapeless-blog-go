package main

import (
	"flag"
	"fmt"
	"os"
)

type ServerCache struct {
	PostIDs    []int
	Categories []string
	IDCounter  int
	TokenHash  []byte
}

type application struct {
	basedir string
	cache   ServerCache
	port    int
}

func main() {
	var basedir string
	var port int
	var createUser bool
	var setup bool
	flag.StringVar(&basedir, "basedir", os.Getenv("HOME")+"/slblog", "Base directory of data")
	flag.IntVar(&port, "port", 9398, "Port of shapeless blog")
	flag.BoolVar(&createUser, "create", false, "Create new user interactively.")
	flag.BoolVar(&setup, "setup", false, "Create necessary files.")

	flag.Parse()

	app := &application{
		basedir: basedir,
		port:    port,
	}

	if setup {
		app.setup()
		return
	}

	app.loadCache()

	if createUser {
		fmt.Printf("Creating new user.\n")
		fmt.Print("Username: ")
		var username string
		fmt.Scanln(&username)
		fmt.Print("Password: ")
		var password string
		fmt.Scanln(&password)

		app.createUser(username, password)
		return
	}

	app.serve()
}

func (app *application) setup() {
	os.Mkdir(app.basedir, os.ModePerm)
	os.Mkdir(app.postsDir(), os.ModePerm)
	os.Mkdir(app.categoryDir(), os.ModePerm)
	os.Mkdir(app.contentDir(), os.ModePerm)
	os.Mkdir(app.metadataDir(), os.ModePerm)
	os.Create(app.postsPath())
	os.Create(app.categoriesPath())
	os.Mkdir(app.templatesDir(), os.ModePerm)
	os.Mkdir(app.usersDir(), os.ModePerm)
	os.Create(app.showBlogPostTemplatePath())
	os.Create(app.showHomePageTemplatePath())
	os.Create(app.showCategoryTemplatePath())
	os.Create(app.showCategoriesIndexTemplatePath())
	f, _ := os.Create(app.IDCounterPath())
	f.WriteString("0\n")
	f.Close()
}
