package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (app *Application) getOnemovie(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))

	if err != nil {
		app.logger.Print(errors.New("Invalid parameter"))
		app.errorJson(w, err)
		return
	}
	movie, err := app.models.DB.Get(id)

	err = app.WriteJSON(w, http.StatusOK, movie, "movie")

	if err != nil {
		app.logger.Print(errors.New("Invalid parameter"))
		app.errorJson(w, err)
		return
	}
}

func (app *Application) getAllmovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.models.DB.GetAll()
	if err != nil {
		app.logger.Print(errors.New("Invalid parameter"))
		app.errorJson(w, err)
		return
	}
	err = app.WriteJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.logger.Print(errors.New("Invalid parameter"))
		app.errorJson(w, err)
		return
	}
}
func (app *Application) GetAllgenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.models.DB.GenresAll()
	if err != nil {
		app.logger.Print(errors.New("Invalid parameter"))
		app.errorJson(w, err)
		return
	}
	err = app.WriteJSON(w, http.StatusOK, genres, "genres")
	if err != nil {
		app.logger.Print(errors.New("Invalid parameter"))
		app.errorJson(w, err)
		return
	}
}
func (app *Application) GetAllmoviesByGenres(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())

	genreID, err := strconv.Atoi(params.ByName("genre_id"))
	if err != nil {
		app.logger.Print(errors.New("Invalid parameter"))
		app.errorJson(w, err)
		return
	}
	movies, err := app.models.DB.GetAll(genreID)
	if err != nil {
		app.logger.Print(errors.New("Invalid parameter"))
		app.errorJson(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, movies, "movies")
	if err != nil {
		app.logger.Print(errors.New("Invalid parameter"))
		app.errorJson(w, err)
		return
	}
}

//Stub functions thath we will remember when we return!!!

type MoviePayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Year        string `json:"year"`
	ReleaseDate string `json:"release_date"`
	Runtime     string `json:"runtime"`
	Rating      string `json:"rating"`
	MPAARating  string `json:"mpaa_rating"`
}

func (app *Application) editmovie(w http.ResponseWriter, r *http.Request) {
	var payload MoviePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	var movie models.Movie
	if payload.ID != "0" {
		id, _ := strconv.Atoi(payload.ID)
		m, _ := app.models.DB.Get(id)
		movie = *m
		movie.UpdatedAt = time.Now()
	}
	movie.ID, _ = strconv.Atoi(payload.ID)
	movie.Title = payload.Title
	movie.Description = payload.Description
	movie.ReleaseDate, err = time.Parse("2006-01-02", payload.ReleaseDate)
	if err != nil {
		log.Fatal(err)
	}
	movie.Year = movie.ReleaseDate.Year()
	movie.Runtime, err = strconv.Atoi(payload.Runtime)
	if err != nil {
		log.Fatal(err)
	}
	movie.Rating, err = strconv.Atoi(payload.Rating)
	if err != nil {
		log.Fatal(err)
	}
	movie.MPAARating = payload.MPAARating
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()
	//If the movie is new then do the insertion!!!
	if movie.ID == 0 {
		err = app.models.DB.InsertMovie(movie)
		if err != nil {
			app.errorJson(w, err)
			return
		}
	} else {
		err = app.models.DB.UpdateMovie(movie)
		if err != nil {
			app.errorJson(w, err)
			return
		}
	}

	ok := jsonResponse{
		OK: true,
	}
	err = app.WriteJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJson(w, err)
		return
	}

}

func (app *Application) DeleteMovie(w http.ResponseWriter, r *http.Request) {

	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	err = app.models.DB.DeleteMovie(id)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	ok := jsonResponse{
		OK: true,
	}
	err = app.WriteJSON(w, http.StatusOK, ok, "response")
	if err != nil {
		app.errorJson(w, err)
		return
	}

}

func (app *Application) SearchMovies(w http.ResponseWriter, r *http.Request) {

}
