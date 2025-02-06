package routes

import (
	"cinemaGo/backend/api/handlers"
	"cinemaGo/backend/api/middlewares"

	"github.com/gin-gonic/gin"
)

type ServeHandlersWrapper struct {
	*handlers.MoviesHandler
	*handlers.UsersHandler
	*handlers.BookingHandler
}

func Router(h *ServeHandlersWrapper) *gin.Engine {

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/home", h.MainPage)
		v1.GET("/explore/movies", h.ExploreAllShows)
		v1.GET("/movies/:movieName/MV:movieID", h.Movie)
		v1.GET("/person/:personName/:actorCrewID", h.ActorCrew)

		v1.POST("/user/signup", h.SignUp)
		v1.POST("/user/login", h.Login)

		v1.GET("/my-profile/edit", middlewares.UserAuthorizationJWT(), h.UserProfile)
		v1.PUT("/my-profile/edit", middlewares.UserAuthorizationJWT(), h.UpdateUserProfile)
		v1.POST("/my-profile/logout", middlewares.UserAuthorizationJWT(), h.Logout)

		v1.GET("/buytickets/movie/:showID/show-times", h.MovieShowTimes)
		v1.GET("/buytickets/movie/:showID/available-seats", h.ShowSeats)

		v1.POST("/buytickets/payment", middlewares.UserAuthorizationJWT(), h.BookSeats)
	}

	return router
}
