package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("starting slblog at port %d\n", app.port)
	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
