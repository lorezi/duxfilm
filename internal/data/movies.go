package data

import "time"

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
	Duration int32    `json:"duration,omitempty"`
	Genres   []string `json:"genres,omitempty"`
	Version  int32    `json:"version"`
}

type MovieRequest struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Year      int32     `json:"year"`
	Duration  int32     `json:"duration"`
	Genres    []string  `json:"genres"`
	Version   int32     `json:"version"`
	CreatedAt time.Time `json:"created_at"`
}
