package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Daniel-Sogbey/hello-world/internal/config"
	"github.com/Daniel-Sogbey/hello-world/internal/handlers"
	"github.com/Daniel-Sogbey/hello-world/internal/models"
	"github.com/Daniel-Sogbey/hello-world/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	err := run()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Server started and listening on port %s", portNumber)
	// http.ListenAndServe(portNumber, nil)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}

func run() error {
	//change this to true when in production
	app.InProduction = false

	//
	gob.Register(models.Reservation{})

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("cannot create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app) //Gives NewRepo the memory address of the app config

	handlers.NewHandlers(repo)

	render.NewTemplates(&app) // Gives NewTemplates the memory address of the app config

	return nil
}
