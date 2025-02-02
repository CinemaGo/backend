package handlers

import (
	"cinemaGo/backend/api/helpers"
	"cinemaGo/backend/internal/services"
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

	c.JSON(http.StatusOK, gin.H{
		"carouselImages": carouselImages,
	})
}
