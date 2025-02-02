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

	movies, err := service.movie.FetchAllMovies()
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"carouselImages": carouselImages,
		"movies":         movies,
	})
}

func (service *MoviesHandler) ExploreAllMovies(c *gin.Context) {
	movies, err := service.movie.FetchAllMovies()
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movies": movies,
	})
}

func (service *MoviesHandler) Movie(c *gin.Context) {
	movieID, err := helpers.GetParameterFromURL(c, "movieID", "Invalid movie ID provided.")
	if err != nil {
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	movies, err := service.movie.FetchAMovie(movieID)
	if err != nil {
		if errors.Is(err, services.ErrMovieNotFoundByID) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("movie with the given ID %v not found", movieID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	actorsCrews, err := service.movie.FetchAllActorsCrewsByMovieID(movieID)
	if err != nil {
		if errors.Is(err, services.ErrMovieNotFoundByID) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("movie with the given ID %v not found", movieID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movies":      movies,
		"actorsCrews": actorsCrews,
	})
}
