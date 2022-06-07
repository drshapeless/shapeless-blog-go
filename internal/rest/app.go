package rest

import (
	"log"
	"sync"

	"github.com/drshapeless/shapeless-blog/internal/data"
)

type Application struct {
	Port      int
	DBPath    string
	Models    data.Models
	InfoLog   *log.Logger
	ErrorLog  *log.Logger
	Wg        sync.WaitGroup
	Version   string
	BuildTime string
	Secret    string
}
