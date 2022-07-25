package main

import (
	"bookings/pkg/config"
	"bookings/pkg/handlers"
	"bookings/pkg/render"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
)

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {

	appConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction

	appConfig.Session = session

	tc, err := render.CreateTemplateCache() //Get all of our template files, store them in a template map so they do not need to be re-rendered
	if err != nil {
		log.Fatal("cannot create template cache")
	}

	appConfig.TemplateCache = tc //Assign the template cach to our app config so they can be accessed elsewhere or use
	appConfig.UseCache = true

	repo := handlers.NewRepo(&appConfig) //Give our 'repo' access to our app config. This means it can associate our global settings
	handlers.NewHandlers(repo)
	render.NewTemplates(&appConfig) //Give the render class access to the app config so it can use the cached templates

	srv := &http.Server{ //Configure the server
		Handler: routes(&appConfig),
	}

	err = srv.ListenAndServe() //Start the server
	log.Fatal(err)

}
