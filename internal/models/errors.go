package models

import "errors"

var ErrMovieNotFoundByID = errors.New("models: movie with the given ID not found")
var ErrActorCrewNotFoundByID = errors.New("models: actor or crew with the given ID not found")

var ErrDuplicatedEmail = errors.New("models: email already exists")
var ErrUserNotFound = errors.New("models: user not found")

var ErrShowNotFound = errors.New("models: show not found by given id")
var ErrStartTimeNotFound = errors.New("models: start time not found")
var ErrShowSeatNotFound = errors.New("models: show seat not found")
