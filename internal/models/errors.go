package models

import "errors"

var ErrMovieNotFoundByID = errors.New("movie with the given ID not found")
var ErrActorCrewNotFoundByID = errors.New("actor or crew with the given ID not found")


var ErrDuplicatedEmail = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")
