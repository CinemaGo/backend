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
	*handlers.AdminHandler
}

func Router(h *ServeHandlersWrapper) *gin.Engine {

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.GET("/home", h.MainPage)
		v1.GET("/explore/movies", h.ExploreAllShows)
		v1.GET("/movies/:movieName/:showID", h.MovieShow)
		v1.GET("/person/:personName/:actorCrewID", h.ActorCrew)

		v1.POST("/user/signup", h.SignUp)
		v1.POST("/user/login", h.Login)

		v1.GET("/my-profile", middlewares.UserAuthorizationJWT(), h.UserProfile)
		v1.PUT("/my-profile/edit", middlewares.UserAuthorizationJWT(), h.UpdateUserProfile)
		v1.POST("/my-profile/logout", middlewares.UserAuthorizationJWT(), h.Logout)

		v1.GET("/buytickets/movie/:showID/show-times", h.MovieShowTimes)
		v1.GET("/buytickets/movie/:showID/available-seats", h.ShowSeats)

		v1.POST("/buytickets/payment", middlewares.UserAuthorizationJWT(), h.BookSeats)

		v1.GET("/admin/carousel-image/all", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.CarouselImagesAdmin)
		v1.POST("/admin/carousel-image/new", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.NewCarouselImageAdmin)
		v1.PUT("/admin/carousel-image/edit", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.EditCarouselImageAdmin)
		v1.DELETE("/admin/carousel-image/delete", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.DeleteCarouselImageAdmin)

		v1.GET("/admin/movie/all", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.AllMoviesAdmin)
		v1.POST("/admin/movie/new", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.NewMovieAdmin)
		v1.GET("/admin/movie/:movieID", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.MovieInfoAdmin)
		v1.PUT("/admin/movie/edit", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.EditMovieAdmin)
		v1.DELETE("/admin/movie/delete", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.DeleteMovieAdmin)

		v1.POST("/admin/movie/actor-crew/new", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.NewActorCrewAdmin)
		v1.GET("/admin/movie/actor-crew/:actorCrewID", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.AnActorCrewAdmin)
		v1.PUT("/admin/movie/actor-crew/edit", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.EditActorCrewAdmin)
		v1.DELETE("/admin/movie/actor-crew/delete", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.DeleteActorCrewAdmin)

		v1.POST("/admin/cinema-hall/new", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.NewCinemaHallAdmin)
		v1.GET("/admin/cinema-hall/all", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.AllCinemaHallsAdmin)
		v1.GET("/admin/cinema-hall/:cinemaHallID", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.CinemaHallAdmin)
		v1.DELETE("/admin/cinema-hall/delete", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.DeleteCinemaHallAdmin)

		v1.POST("/admin/cinema-hall-seat/new", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.NewCinemaHallSeatAdmin)
		v1.DELETE("/admin/cinema-hall-seat/delete", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.DeleteCinemaHallSeatAdmin)

		v1.GET("/admin/show/all", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.AllShowsAdmin)
		v1.POST("/admin/show/new", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.NewShowAdmin)
		v1.PUT("/admin/show/edit", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.EditShowAdmin)
		v1.DELETE("/admin/show/delete", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.DeleteShowAdmin)

		v1.GET("/admin/show-seats/:showID", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.AllShowSeatsAdmin)
		v1.PUT("/admin/show-seat-price/edit", middlewares.UserAuthorizationJWT(), middlewares.AdminRoleRequired(), h.EditShowSeatPriceAdmin)

	}

	return router
}
