package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// This is for webpage.
	router.GET("/", app.showHomePageHandler)
	router.GET("/post/:id", app.showBlogpostHandler)
	router.GET("/category/", app.showCategoriesIndexHandler)
	router.GET("/category/:category", app.showCategoryHandler)

	// The following is for api.
	router.POST("/api/v1/post/", app.requireAuthentication(app.createBlogpostHandler))
	router.GET("/api/v1/post/:id", app.requireAuthentication(app.getBlogpostHandler))
	router.PATCH("/api/v1/post/:id", app.requireAuthentication(app.updateBlogpostHandler))
	router.DELETE("/api/v1/post/:id", app.requireAuthentication(app.deleteBlogpostHandler))

	// Authentication.
	router.POST("/api/v1/authentication/", app.authenticationHandler)

	return router
}
