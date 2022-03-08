package main

import "fmt"

func (app *application) databasePath() string {
	return app.config.dir + "/shapeless-blog.db"
}

func (app *application) templatePath(name string) string {
	return app.config.dir + "/templates/" + name + ".tmpl"
}

func (app *application) postPath(filename string) string {
	return fmt.Sprintf("%s/posts/%s.html", app.config.dir, filename)
}
