package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Movies interface {
		Insert(movie *Movie) error
		Get(id int64) (*Movie, error)
		Update(movie *Movie) error
		Delete(id int64) error
	}
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}
}

// Create a helper function which returns a Models instance containing the rock models only.

/*
COME BACK TO THIS LINE OF CODE
*/
// func NewMockModels() Models {
// 	return Models{
// 		Movies: MockMovieModel{},
// 	}
// }
