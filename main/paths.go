package main

import "fmt"

func (app *application) postsDir() string {
	return app.basedir + "/posts"
}

func (app *application) postsPath() string {
	return app.postsDir() + "/posts.txt"
}

func (app *application) IDCounterPath() string {
	return app.postsDir() + "/idcounter.txt"
}

func (app *application) categoriesPath() string {
	return app.postsDir() + "/categories.txt"
}

func (app *application) metadataDir() string {
	return app.postsDir() + "/metadata"
}

func (app *application) metadataPath(id int) string {
	// By prepending some zeros in the filename, it makes dired to have
	// a better order of displaying.
	return app.metadataDir() + fmt.Sprintf("/%08d.json", id)
}

func (app *application) contentDir() string {
	return app.postsDir() + "/content"
}

func (app *application) contentPath(id int) string {
	return app.contentDir() + fmt.Sprintf("/%08d.txt", id)
}

func (app *application) categoryDir() string {
	return app.postsDir() + "/category"
}

func (app *application) categoryPath(name string) string {
	return app.categoryDir() + "/" + name + ".txt"
}

func (app *application) templatesDir() string {
	return app.basedir + "/templates"
}

func (app *application) showBlogPostTemplatePath() string {
	return app.templatesDir() + "/showpost.tmpl"
}

func (app *application) showHomePageTemplatePath() string {
	return app.templatesDir() + "/home.tmpl"
}

func (app *application) showCategoriesIndexTemplatePath() string {
	return app.templatesDir() + "/categories_index.tmpl"
}

func (app *application) showCategoryTemplatePath() string {
	return app.templatesDir() + "/category.tmpl"
}

func (app *application) usersDir() string {
	return app.basedir + "/users"
}

func (app *application) userPath(user string) string {
	return app.usersDir() + "/" + user + ".txt"
}

func (app *application) tokenPath() string {
	return app.basedir + "/token.txt"
}
