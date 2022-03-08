package main

import (
	"html/template"
	"log"
	"net/http"
)

func (app *application) errorResponse(w http.ResponseWriter, status int, message string) {
	tmpl, err := template.ParseFiles(app.templatePath("error"))
	if err != nil {
		log.Fatal("really fucked up")
		return
	}
	w.WriteHeader(status)

	var body struct {
		Message string
		Status  int
	}
	body.Message = message
	body.Status = status

	tmpl.Execute(w, body)
}

func (app *application) serverErrorResponse(w http.ResponseWriter, err error) {
	log.Print(err)

	message := "the requested resource could not be found"
	app.errorResponse(w, http.StatusNotFound, message)
}
