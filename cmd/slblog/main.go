package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/drshapeless/shapeless-blog/internal/data"

	_ "github.com/mattn/go-sqlite3"
)

type config struct {
	port int
	path string
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

}

func openDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	return db, nil
}
