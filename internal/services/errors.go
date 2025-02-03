package services

import "errors"

var ErrMovieNotFoundByID = errors.New("movie with the given ID not found")
var ErrActorCrewNotFoundByID = errors.New("actor or crew with the given ID not found")
var ErrNoCachedDataFound = errors.New("no data found in the Redis cache")
