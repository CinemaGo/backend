package services

import "errors"

var ErrMovieNotFoundByID = errors.New("movie with the given ID not found")
var ErrActorCrewNotFoundByID = errors.New("actor or crew with the given ID not found")
var ErrNoCachedDataFound = errors.New("no data found in the Redis cache")

var ErrDuplicatedEmail = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrUserInvalidCredentials = errors.New("invalid credentials")

var ErrShowNotFound = errors.New("models: show not found by given id")
var ErrStartTimeNotFound = errors.New("models: start time not found")
var ErrShowSeatNotFound = errors.New("models: show seat not found")

var ErrShowSeatHasSelected = errors.New("show seat has just selected or booked")
var ErrTooManySeats = errors.New("too many seats selected")
