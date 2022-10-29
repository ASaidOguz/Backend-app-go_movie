package main

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) wrap(next http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := context.WithValue(r.Context(), "params", ps)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (app *application) routes() http.Handler {
	router := httprouter.New()
	secure := alice.New(app.checkToken)
	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)

	router.HandlerFunc((http.MethodPost), "/v1/graphql/list", app.moviesGraphql)

	router.HandlerFunc(http.MethodPost, "/v1/signin", app.Signin)
	router.HandlerFunc(http.MethodGet, "/v1/movie/:id", app.getOnemovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.getAllmovies)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:genre_id", app.GetAllmoviesByGenres)

	router.HandlerFunc(http.MethodGet, "/v1/genres", app.GetAllgenres)

	router.POST("/v1/admin/editmovie", app.wrap(secure.ThenFunc(app.editmovie)))
	//router.HandlerFunc(http.MethodPost, "/v1/admin/editmovie", app.editmovie)
	router.GET("/v1/admin/deletemovie/:id", app.wrap(secure.ThenFunc(app.DeleteMovie)))
	//router.HandlerFunc(http.MethodGet, "/v1/admin/deletemovie/:id", app.DeleteMovie)
	return app.enableCORS(router)
}
