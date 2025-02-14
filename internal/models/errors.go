package models

import "errors"

var ErrMovieNotFoundByID = errors.New("models: movie with the given ID not found")
var ErrActorCrewNotFoundByID = errors.New("models: actor or crew with the given ID not found")

var ErrDuplicatedEmail = errors.New("models: email already exists")
var ErrUserNotFound = errors.New("models: user not found")

var ErrShowNotFound = errors.New("models: show not found by given id")
var ErrStartTimeNotFound = errors.New("models: start time not found")
var ErrShowSeatNotFound = errors.New("models: show seat not found")

var ErrAdminPageCarouselImagesNotFound = errors.New("models: Admin Page, Carousel Images Not Found")
var ErrAdminPageMovieNotFound = errors.New("models: Admin Page, Movie Not Found")
var ErrActorCrewNotFound = errors.New("models: Admin page, actorCrew with ID not found")
var ErrDuplicatedCinemaHall = errors.New("models: Admin page, cinema hall with this name and type already exists")
var ErrCinemaHallNotFound = errors.New("models: Admin page, cinema hall not found")
var ErrCinemaSeatNotFound = errors.New("models: Admin page, cinema seat not found")
var ErrCinemaSeatAlreadyExists = errors.New("models: Admin page, cinema seat with hall_id, seat_row, seat_number already exists")
var ErrShowAlreadyExists = errors.New("models: Admin page, a show already exists at the given hall, date, and time")