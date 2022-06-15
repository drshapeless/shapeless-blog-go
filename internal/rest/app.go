package rest

import (
	"html/template"
	"log"
	"sync"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

type Application struct {
	Port          int
	DBPath        string
	Models        data.Models
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	Wg            sync.WaitGroup
	Version       string
	BuildTime     string
	Secret        string
	TemplateCache map[string]*template.Template
}
