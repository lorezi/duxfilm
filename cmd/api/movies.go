package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/lorezi/duxfilm/internal/data"
	"github.com/lorezi/duxfilm/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title    string        `json:"title"`
		Year     int32         `json:"year"`
		Duration data.Duration `json:"duration"`
		Genres   []string      `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// validation
	v := validator.New()

	v.Check(input.Title != "", "title", "must be provided")
	v.Check(len(input.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(input.Year != 0, "year", "must be provided")
	v.Check(input.Year >= 1888, "year", "must be greater than 1888")
	v.Check(input.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(input.Duration != 0, "duration", "must be provided")
	v.Check(input.Duration > 0, "duration", "must be a positive integer")

	v.Check(input.Genres != nil, "genres", "must be provided")
	v.Check(len(input.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(input.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
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
