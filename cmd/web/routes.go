package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) routes() http.Handler {
	// We are gonna use mux router
	mux := chi.NewRouter() // Call the function to get a router instance

	// Set up middleware
	mux.Use(middleware.Recoverer)
	mux.Use(app.SessionLoad)

	// Define application routes
	mux.Get("/", app.HomePage)

	mux.Get("/login", app.LogInPage)
	mux.Post("/login", app.postLoginPage)
	mux.Get("/logout", app.LogOutPage)
	mux.Get("/register", app.registerPage)
	mux.Post("/register", app.postRegisterPage)
	mux.Get("/activate-account", app.activateAccount)

	return mux
}
