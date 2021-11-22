package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies     MovieModel
	Tokens     TokenModel
	User       UserModel
	Permission PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:     MovieModel{DB: db},
		Tokens:     TokenModel{DB: db},
		User:       UserModel{DB: db},
		Permission: PermissionModel{DB: db},
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
