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

var ErrAdminPageCarouselImagesNotFound = errors.New("admin Page, Carousel Images Not Found")
var ErrAdminPageMovieNotFound = errors.New("admin Page, Movie Not Found")
var ErrActorCrewNotFound = errors.New("admin page, actorCrew with ID not found")
var ErrDuplicatedCinemaHall = errors.New("admin page, cinema hall with this name and type already exists")
var ErrCinemaHallNotFound = errors.New("admin page, cinema hall not found")
var ErrCinemaSeatNotFound = errors.New("admin page, cinema seat not found")
var ErrCinemaSeatAlreadyExists = errors.New("admin page, cinema seat with hall_id, seat_row, seat_number already exists")
var ErrShowAlreadyExists = errors.New("admin page, a show already exists at the given hall, date, and time")
