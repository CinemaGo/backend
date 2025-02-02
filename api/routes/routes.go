package routes

import (
	"cinemaGo/backend/api/handlers"

	"github.com/gin-gonic/gin"
)

type ServeHandlersWrapper struct {
	*handlers.MoviesHandler
}

func Router(h *ServeHandlersWrapper) *gin.Engine {

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/", h.MainPage)
	}

	return router
}
