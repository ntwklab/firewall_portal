package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/ntwklab/firewall_portal/internal/config"
	"github.com/ntwklab/firewall_portal/internal/handlers"
	"github.com/ntwklab/firewall_portal/internal/helpers"
	"github.com/ntwklab/firewall_portal/internal/models"
	"github.com/ntwklab/firewall_portal/internal/render"
)

const portNumber = ":8010"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

// main is the main
func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Starting application on port " + portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)

}

func run() error {

	// What am i going to put in the session
	gob.Register(models.CreateRule{})

	// change this to true when in production
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFT\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

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

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	render.NewTemplates(&app)
	helpers.NewHelpers(&app)

	return nil
}
