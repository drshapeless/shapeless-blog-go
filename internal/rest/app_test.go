package rest_test

import (
	"log"
	"os"
	"sync"

	"github.com/drshapeless/shapeless-blog/internal/data"
	"github.com/drshapeless/shapeless-blog/internal/rest"
)

func TestApplication() *rest.Application {
	p := "shapeless-blog.db"
	db, err := data.OpenDB(p)
	if err != nil {
		return nil
	}

	app := rest.Application{
		Port:     9398,
		DBPath:   "shapeless-blog.db",
		Models:   data.NewModels(db),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Wg:       sync.WaitGroup{},
	}

	return &app
}
