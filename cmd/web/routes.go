package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ntwklab/firewall_portal/internal/config"
	"github.com/ntwklab/firewall_portal/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	mux.Get("/create-rule", handlers.Repo.CreateRule)
	mux.Post("/create-rule", handlers.Repo.PostCreateRule)
	mux.Get("/create-rule-summary", handlers.Repo.CreateRuleSummary)

	mux.Post("/check-duplicate", handlers.Repo.CheckDuplicate)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
