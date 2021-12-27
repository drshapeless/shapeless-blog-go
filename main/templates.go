package main

import "text/template"

func (app *application) showBlogpostTemplate() *template.Template {
	tmpl, err := template.ParseFiles(app.showBlogPostTemplatePath())
	if err != nil {
		// The template must be there.
		panic(err)
	}
	return tmpl
}

func (app *application) showHomePageTemplate() *template.Template {
	tmpl, err := template.ParseFiles(app.showHomePageTemplatePath())
	if err != nil {
		panic(err)
	}
	return tmpl
}

func (app *application) showCategoriesIndexTemplate() *template.Template {
	tmpl, err := template.ParseFiles(app.showCategoriesIndexTemplatePath())
	if err != nil {
		panic(err)
	}

	return tmpl
}

func (app *application) showCategoryTemplate() *template.Template {
	tmpl, err := template.ParseFiles(app.showCategoryTemplatePath())
	if err != nil {
		panic(err)
	}

	return tmpl
}
