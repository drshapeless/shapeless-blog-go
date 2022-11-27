package rest

import (
	"html/template"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

func (app *Application) templateHelper(name string) *template.Template {
	ts, err := app.Models.Templates.GetByName(name)
	if err != nil {
		return nil
	}

	tmpl, err := template.New(name).Parse(ts.Content)
	if err != nil {
		return nil
	}

	return tmpl
}

// Make sure to call this function after app.Models is populated.
func (app *Application) RefreshTemplateCache() map[string]*template.Template {
	cache := map[string]*template.Template{}

	cache["home"] = app.templateHelper("home")
	cache["post"] = app.templateHelper("post")
	cache["tag"] = app.templateHelper("tag")
	cache["list-tags"] = app.templateHelper("list-tags")

	return cache
}

func (app *Application) updateTemplateCache(t *data.Template) {
	tmpl, err := template.New(t.Name).Parse(t.Content)
	if err != nil {
		// Fuck that, whatever.
		return
	}

	app.TemplateCache[t.Name] = tmpl
}
