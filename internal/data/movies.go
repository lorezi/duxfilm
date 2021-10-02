package data

import (
	"time"

	"github.com/lorezi/duxfilm/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Year      int32     `json:"year"`
	Duration  int32     `json:"duration"`
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}

type MovieResponse struct {
	ID       int64    `json:"id"`
	Title    string   `json:"title"`
	Year     int32    `json:"year,omitempty"`
	Duration Duration `json:"duration,omitempty"`
	Genres   []string `json:"genres,omitempty"`
	Version  int32    `json:"version"`
}

type MovieRequest struct {
	ID       int64    `json:"id"`
	Title    string   `json:"title"`
	Year     int32    `json:"year"`
	Duration int32    `json:"duration"`
	Genres   []string `json:"genres"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Duration != 0, "duration", "must be provided")
	v.Check(movie.Duration > 0, "duration", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
