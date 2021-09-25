package main

import (
	"fmt"
	"net/http"

	"github.com/lorezi/duxfilm/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new movie")
}

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
