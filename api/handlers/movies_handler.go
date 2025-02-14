package handlers

import (
	"cinemaGo/backend/api/helpers"
	"cinemaGo/backend/internal/services"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MoviesHandler struct {
	movie services.MoviesServiceInterface
}

func NewMoviesHandler(service services.MoviesServiceInterface) *MoviesHandler {
	return &MoviesHandler{movie: service}
}

func (service *MoviesHandler) MainPage(c *gin.Context) {

	carouselImages, err := service.movie.FetchAllCaruselImages()
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	showsMovie, err := service.movie.FetchAllShowsMovie()

	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"carouselImages": carouselImages,
		"showsMovie":     showsMovie,
	})
}

func (service *MoviesHandler) ExploreAllShows(c *gin.Context) {
	showsMovie, err := service.movie.FetchAllShowsMovie()
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"showsMovie": showsMovie,
	})
}

func (service *MoviesHandler) MovieShow(c *gin.Context) {
	showID, err := helpers.GetParameterFromURL(c, "showID", "invalid showID ID provided.")
	if err != nil {
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	aShowMovie, err := service.movie.FetchAShowMovie(showID)
	if err != nil {
		if errors.Is(err, services.ErrMovieNotFoundByID) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("showID with the given ID %v not found", showID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	actorsCrews, err := service.movie.FetchAllActorsCrewsByMovieID(aShowMovie.MovieID)
	if err != nil {
		if errors.Is(err, services.ErrMovieNotFoundByID) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("showID with the given ID %v not found", showID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aShowMovie":  aShowMovie,
		"actorsCrews": actorsCrews,
	})
}

func (service *MoviesHandler) ActorCrew(c *gin.Context) {
	actorCrewID, err := helpers.GetParameterFromURL(c, "actorCrewID", "invalid person ID provided.")
	if err != nil {
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	actorCrew, err := service.movie.FetchActorCrewInfo(actorCrewID)
	if err != nil {
		if errors.Is(err, services.ErrActorCrewNotFoundByID) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("persone with the gived ID %v not found", actorCrewID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	movies, err := service.movie.FetchMoviesByActorCrewID(actorCrewID)
	if err != nil {
		if errors.Is(err, services.ErrActorCrewNotFoundByID) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("persone with the gived ID %v not found", actorCrewID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"person": actorCrew,
		"movies": movies,
	})
}
