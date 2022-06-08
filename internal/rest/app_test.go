package rest

import (
	"log"
	"os"
	"sync"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

func ExampleApplication() *Application {
	p := os.Getenv("SHAPELESS_BLOG_DB_PATH")
	db, err := data.OpenDB(p)
	if err != nil {
		return nil
	}

	app := Application{
		Port:     9398,
		DBPath:   p,
		Models:   data.NewModels(db),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		Wg:       sync.WaitGroup{},
		Secret:   "testsecret",
	}

	return &app
}
