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

	// Define application routes
	mux.Get("/", app.HomePage)

	return mux
}
