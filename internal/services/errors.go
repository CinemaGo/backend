package services

import "errors"

var ErrMovieNotFoundByID = errors.New("movie with the given ID not found")
var ErrActorCrewNotFoundByID = errors.New("actor or crew with the given ID not found")
var ErrNoCachedDataFound = errors.New("no data found in the Redis cache")

var ErrDuplicatedEmail = errors.New("email already exists")
var ErrUserNotFound = errors.New("user not found")
var ErrUserInvalidCredentials = errors.New("invalid credentials")

var ErrShowNotFound = errors.New("show not found by given id")
var ErrStartTimeNotFound = errors.New("start time not found")
var ErrShowSeatNotFound = errors.New("show seat not found")

var ErrShowSeatHasSelected = errors.New("show seat has just selected or booked")
var ErrTooManySeats = errors.New("too many seats selected")

var ErrAdminPageCarouselImagesNotFound = errors.New("Admin Page Carousel Images Not Found")
var ErrAdminPageMovieNotFound = errors.New("Admin Page Movie Not Found")
var ErrActorCrewNotFound = errors.New("Admin page actorCrew with ID not found")
var ErrCinemaHallNotFound = errors.New("Admin page cinema hall not found")
var ErrCinemaSeatNotFound = errors.New("Admin page cinema seat not found")
var ErrCinemaSeatAlreadyExists = errors.New("cinema seat with hall_id, seat_row, seat_number already exists")
