package main

import (
	"flag"
	"fmt"
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
	flag.StringVar(&app.DBPath, "path", os.Getenv("SHAPELESS_BLOG_DB_PATH"), "shapeless-blog database path")
	flag.IntVar(&app.Port, "port", 9398, "shapeless-blog port")
	flag.StringVar(&app.Secret, "secret", os.Getenv("SHAPELESS_BLOG_SECRET"), "shapeless-blog secret")
	migrate := flag.Bool("migrate", false, "migrate database")
	reset := flag.Bool("reset", false, "reset database")
	displayVersion := flag.Bool("version", false, "Display version and exit")
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

	if *displayVersion {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build time: %s\n", buildTime)
		os.Exit(0)
	}

	if *migrate {
		app.InfoLog.Println("Migrating database...")
		err = data.Migrate(db)
		if err != nil {
			app.ErrorLog.Fatalln(err)
		}
		app.InfoLog.Println("Migration done.")
		return
	}

	if *reset {
		app.InfoLog.Println("Resetting database...")
		err = data.ResetDB(db)
		if err != nil {
			app.ErrorLog.Fatalln(err)
		}
		app.InfoLog.Println("Reset done.")
		return
	}

	app.Models = data.NewModels(db)

	app.InfoLog.Println("database connection pool established")

	app.TemplateCache = app.RefreshTemplateCache()

	app.InfoLog.Println("template cache loaded")

	err = app.Serve()
	if err != nil {
		app.ErrorLog.Fatalln(err)
	}
}
