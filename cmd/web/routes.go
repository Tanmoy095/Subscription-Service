package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *Config) routes() http.Handler {
	//we are gonna use mux router
	mux := chi.NewRouter
	//setUp Middleware
	mux.use(middleware.Recoverer)
	//define application routes

	mux.Get("/", app.HomePage)
	return mux

}
