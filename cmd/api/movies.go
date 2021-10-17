package main

import (
	"fmt"
	"net/http"

	"github.com/lorezi/duxfilm/internal/data"
	"github.com/lorezi/duxfilm/internal/validator"
)

func (app *application) getMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.getParamID(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	movie := data.MovieResponse{
		ID:       id,
		Title:    "Godfather",
		Duration: 202,
		Genres:   []string{"drama", "action"},
		Version:  1,
	}

	// movieres := data.Envelope{
	// 	Movie: movie,
	// }

	// encode the movie data
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var req data.MovieRequest
	// 1. Read JSON to Object
	err := app.readJSON(w, r, &req)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &data.Movie{
		Title:    req.Title,
		Year:     req.Year,
		Duration: int32(req.Duration),
		Genres:   req.Genres,
	}

	//2. Validate the data
	v := validator.New()

	// 3. Store the data
	if data.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	err = app.models.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"movie": movie}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
