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
		v1.GET("/home", h.MainPage)
		v1.GET("/explore/movies",h.ExploreAllShows)
		v1.GET("/movies/:movieName/MV:movieID",h.Movie)
		v1.GET("/person/:personName/:actorCrewID",h.ActorCrew)
	}

	return router
}
