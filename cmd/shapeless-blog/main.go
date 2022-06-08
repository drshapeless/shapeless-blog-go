package main

import (
	"flag"
	"log"
	"os"

	"github.com/drshapeless/shapeless-blog/internal/data"
	"github.com/drshapeless/shapeless-blog/internal/rest"
)

var (
	version   string
	buildTime string
)

func main() {
	var app rest.Application
	flag.StringVar(&app.DBPath, "dir", os.Getenv("SHAPELESS_BLOG_DB_PATH"), "shapeless-blog database path")
	flag.IntVar(&app.Port, "port", 9398, "shapeless-blog port")
	flag.StringVar(&app.Secret, "secret", os.Getenv("SHAPELESS_BLOG_SECRET"), "shapeless-blog secret")
	migrate := flag.Bool("migrate", false, "migrate database")
	flag.Parse()

	app.Version = version
	app.BuildTime = buildTime

	app.InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := data.OpenDB(app.DBPath)
	if err != nil {
		app.ErrorLog.Fatalln(err)
	}
	defer db.Close()

	if *migrate {
		app.InfoLog.Println("Migrating database")
		err = data.Migrate(db)
		if err != nil {
			app.ErrorLog.Fatalln(err)
		}
		return
	}

	app.InfoLog.Println("database connection pool established")

	err = app.Serve()
	if err != nil {
		app.ErrorLog.Fatalln(err)
	}
}
