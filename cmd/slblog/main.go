package main

import (
	"database/sql"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/drshapeless/shapeless-blog/internal/data"

	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	port int
	path string
}

type application struct {
	config config
	models data.Models
}

func main() {
	var cfg config
	flag.StringVar(&cfg.path, "path", "slblog.db", "Path to sqlite datebase")
	flag.IntVar(&cfg.port, "port", 8393, "slblog port")

	install := flag.Bool("install", false, "Install slblog database")
	flag.Parse()

	if *install {
		db, err := openDB(cfg.path)
		if err != nil {
			log.Fatal(err)
		}
		data.Migrate(db)
		log.Printf("Successfully installed slblog database\n")
		return
	}

	if !exists(cfg.path) {
		log.Printf("Database not found at %s", cfg.path)
		log.Printf("Please use -install to install database first.")
		return
	}

	db, err := openDB(cfg.path)
	if err != nil {
		log.Fatal(err)
	}

	var app application
	app.models = data.NewModels(db)
	app.config = cfg

	app.serve()
}

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
