package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc(http.MethodPost, "/v1/signin", app.signin)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOnemovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllmovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:genre_id", app.GetAllmoviesByGenres)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.GetAllgenres)

	router.HandlerFunc(http.MethodPost, "/v1/admin/editmovie", app.editmovie)
	router.HandlerFunc(http.MethodGet, "/v1/admin/deletemovie/:id", app.DeleteMovie)
	return app.enableCORS(router)
}
