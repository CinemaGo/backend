package handlers

import (
	"cinemaGo/backend/api/helpers"
	"cinemaGo/backend/internal/models"
	"cinemaGo/backend/internal/services"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminCtrl services.AdminServiceInterface
}

func NewAdminHandler(service services.AdminServiceInterface) *AdminHandler {
	return &AdminHandler{adminCtrl: service}
}

func (service *AdminHandler) CarouselImagesAdmin(c *gin.Context) {
	allCarouselImages, err := service.adminCtrl.FetchAllCarouselImages()
	if err != nil {
		if errors.Is(err, services.ErrAdminPageCarouselImagesNotFound) {
			helpers.ClientError(c, http.StatusNotFound, "There is no carousel images yet!")
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"allCarouselImages": allCarouselImages,
	})
}

func (service *AdminHandler) NewCarouselImageAdmin(c *gin.Context) {
	var newCarouselImage NewCarouselImageForm

	if err := c.ShouldBindJSON(&newCarouselImage); err != nil {
		helpers.RespondWithValidationErrors(c, err, newCarouselImage)
		return
	}

	err := service.adminCtrl.CreateNewCarouselImage(newCarouselImage.ImageURL, newCarouselImage.Title, newCarouselImage.Description, newCarouselImage.OrderPriority)
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "New carousel image added successfully!",
	})
}

func (service *AdminHandler) EditCarouselImageAdmin(c *gin.Context) {
	var carouselImage EditCarouselImageForm

	if err := c.ShouldBindJSON(&carouselImage); err != nil {
		helpers.RespondWithValidationErrors(c, err, carouselImage)
		return
	}

	err := service.adminCtrl.UpdateCarouselImages(carouselImage.ImageURL, carouselImage.Title, carouselImage.Description, carouselImage.OrderPriority, carouselImage.CarouselImageID)
	if err != nil {
		if errors.Is(err, services.ErrAdminPageCarouselImagesNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("Carousel image by provided %v ID not found", carouselImage.CarouselImageID))
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Carousel image data updated successfully!",
	})
}

func (service *AdminHandler) DeleteCarouselImageAdmin(c *gin.Context) {
	var carouselImage DeleteCarouselImageForm

	if err := c.ShouldBindJSON(&carouselImage); err != nil {
		helpers.RespondWithValidationErrors(c, err, carouselImage)
		return
	}

	err := service.adminCtrl.DeleteCarouselImages(carouselImage.CarouselImageID)
	if err != nil {
		if errors.Is(err, services.ErrAdminPageCarouselImagesNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("Carousel image by provided %v ID not found", carouselImage.CarouselImageID))
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Carousel image data updated successfully!",
	})
}

func (service *AdminHandler) AllMoviesAdmin(c *gin.Context) {
	allMovies, err := service.adminCtrl.FetchAllMovies()
	if err != nil {
		if errors.Is(err, services.ErrAdminPageMovieNotFound) {
			helpers.ClientError(c, http.StatusNotFound, "There is no movie yet!")
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"allMovies": allMovies,
	})
}

func (service *AdminHandler) NewMovieAdmin(c *gin.Context) {
	var newMovie NewMovieForm

	if err := c.ShouldBindJSON(&newMovie); err != nil {
		helpers.RespondWithValidationErrors(c, err, newMovie)
		return
	}

	err := service.adminCtrl.AddNewMovie(newMovie.Title, newMovie.Description, newMovie.Genre, newMovie.Language, newMovie.TrailerURL, newMovie.PosterURL, newMovie.Rating, newMovie.RatingProvider, newMovie.Duration, newMovie.ReleaseDate, newMovie.AgeLimit)
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "New movie added successfully!",
	})
}

func (service *AdminHandler) MovieInfoAdmin(c *gin.Context) {
	movieID, err := helpers.GetParameterFromURL(c, "movieID", "invalid movie ID provided.")
	if err != nil {
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	movie, err := service.adminCtrl.FetchAMovie(movieID)
	if err != nil {
		if errors.Is(err, services.ErrAdminPageMovieNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("movie provided with %v ID not found", movieID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	allActorsCrew, err := service.adminCtrl.FetchAllActorsCrewByMovieID(movie.MovieID)
	if err != nil {
		if errors.Is(err, services.ErrActorCrewNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("actor/crew not found by provided movie ID %v", movie.MovieID))
			c.JSON(http.StatusOK, gin.H{
				"movie": movie,
			})
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"movie":         movie,
		"allActorsCrew": allActorsCrew,
	})

}

func (service *AdminHandler) EditMovieAdmin(c *gin.Context) {
	var movie EditMovieForm

	if err := c.ShouldBindJSON(&movie); err != nil {
		helpers.RespondWithValidationErrors(c, err, movie)
		return
	}

	err := service.adminCtrl.UpdateMovieInfo(movie.MovieID, movie.Title, movie.Description, movie.Genre, movie.Language, movie.TrailerURL, movie.PosterURL, movie.Rating, movie.RatingProvider, movie.Duration, movie.ReleaseDate, movie.AgeLimit)
	if err != nil {
		if errors.Is(err, services.ErrAdminPageMovieNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("movie provided with %v ID not found", movie.MovieID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie data updated successfully!",
	})
}

func (service *AdminHandler) DeleteMovieAdmin(c *gin.Context) {
	var movie DeleteMovieForm

	if err := c.ShouldBindJSON(&movie); err != nil {
		helpers.RespondWithValidationErrors(c, err, movie)
		return
	}

	err := service.adminCtrl.DeleteMovie(movie.MovieID)
	if err != nil {
		if errors.Is(err, services.ErrAdminPageMovieNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("movie provided with %v ID not found", movie.MovieID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Movie deleted successfully!",
	})
}

func (service *AdminHandler) NewActorCrewAdmin(c *gin.Context) {
	var newActorCrew NewActorCrewForm

	if err := c.ShouldBindJSON(&newActorCrew); err != nil {
		helpers.RespondWithValidationErrors(c, err, newActorCrew)
		return
	}

	err := service.adminCtrl.AddActorsCrew(newActorCrew.FullName, newActorCrew.ImageURL, newActorCrew.Occupation, newActorCrew.RoleDescription, newActorCrew.BornDate, newActorCrew.Birthplace, newActorCrew.About, newActorCrew.IsActor, newActorCrew.MovieID)
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Actor/Crew added successfully!",
	})
}

func (service *AdminHandler) AnActorCrewAdmin(c *gin.Context) {
	actorCrewID, err := helpers.GetParameterFromURL(c, "actorCrewID", "invalid actor/crew ID provided.")
	if err != nil {
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	actorCrew, err := service.adminCtrl.FetchAnActorCrew(actorCrewID)
	if err != nil {
		if errors.Is(err, services.ErrActorCrewNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("actor/crew not found by provided ID %v", actorCrewID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"actorCrew": actorCrew,
	})

}

func (service *AdminHandler) EditActorCrewAdmin(c *gin.Context) {
	var actorCrew EditActorCrewForm

	if err := c.ShouldBindJSON(&actorCrew); err != nil {
		helpers.RespondWithValidationErrors(c, err, actorCrew)
		return
	}

	err := service.adminCtrl.UpdateActorCrewInfo(actorCrew.FullName, actorCrew.ImageURL, actorCrew.Occupation, actorCrew.RoleDescription, actorCrew.BornDate, actorCrew.Birthplace, actorCrew.About, actorCrew.IsActor, actorCrew.ActorCrewID)
	if err != nil {
		if errors.Is(err, models.ErrActorCrewNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("actor/crew not found by provided ID %v", actorCrew.ActorCrewID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Actor/Crew data updated successfully!",
	})
}

func (service *AdminHandler) DeleteActorCrewAdmin(c *gin.Context) {
	var actorCrew DeleteActorCrewForm

	if err := c.ShouldBindJSON(&actorCrew); err != nil {
		helpers.RespondWithValidationErrors(c, err, actorCrew)
		return
	}

	err := service.adminCtrl.DeleteActorCrew(actorCrew.ActorCrewID)
	if err != nil {
		if errors.Is(err, services.ErrActorCrewNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("actor/crew not found by provided ID %v", actorCrew.ActorCrewID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Actor/Crew deleted successfully!",
	})
}

func (service *AdminHandler) AllCinemaHallsAdmin(c *gin.Context) {
	allCinemaHalls, err := service.adminCtrl.FetchAllCinemaHalls()
	if err != nil {
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			helpers.ClientError(c, http.StatusNotFound, "There is no cinema hall yet!")
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"allCinemaHalls": allCinemaHalls,
	})
}

func (service *AdminHandler) NewCinemaHallAdmin(c *gin.Context) {
	var newCinemaHall NewCinemaHallForm

	if err := c.ShouldBindJSON(&newCinemaHall); err != nil {
		helpers.RespondWithValidationErrors(c, err, newCinemaHall)
		return
	}

	err := service.adminCtrl.AddNewCinemaHall(newCinemaHall.HallName, newCinemaHall.HallName, 0)
	if err != nil {
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "New Cinema Hall added successfully!",
	})
}

func (service *AdminHandler) CinemaHallAdmin(c *gin.Context) {
	cinemaHallID, err := helpers.GetParameterFromURL(c, "cinemaHallID", "invalid cinema hall ID provided.")
	if err != nil {
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	cinemaHall, err := service.adminCtrl.FetchCinemaHallInfo(cinemaHallID)
	if err != nil {
		if errors.Is(err, services.ErrCinemaHallNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("cinema hall not found by provided ID %v", cinemaHallID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	allCinemaSeats, err := service.adminCtrl.FetchALLCinemaSeatsByHallID(cinemaHallID)
	if err != nil {
		if errors.Is(err, services.ErrCinemaSeatNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"cinemaHall":     cinemaHall,
				"allCinemaSeats": "These are no seats yet!",
			})
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cinemaHall":     cinemaHall,
		"allCinemaSeats": allCinemaSeats,
	})
}

func (service *AdminHandler) EditCinemaHallAdmin(c *gin.Context) {
	var cinemaHall EditCinemaHallForm

	if err := c.ShouldBindJSON(&cinemaHall); err != nil {
		helpers.RespondWithValidationErrors(c, err, cinemaHall)
		return
	}

	err := service.adminCtrl.UpdateCinemaHallInfo(cinemaHall.HallName, cinemaHall.HallType, 0, cinemaHall.CinemaHallID)
	if err != nil {
		if errors.Is(err, services.ErrCinemaHallNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("cinema hall not found by provided ID %v", cinemaHall.CinemaHallID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cinema Hall data updated successfully",
	})
}

func (service *AdminHandler) DeleteCinemaHallAdmin(c *gin.Context) {
	var cinemaHall DeleteCinemaHallForm

	if err := c.ShouldBindJSON(&cinemaHall); err != nil {
		helpers.RespondWithValidationErrors(c, err, cinemaHall)
		return
	}

	err := service.adminCtrl.DeleteCinemaHall(cinemaHall.CinemaHallID)
	if err != nil {
		if errors.Is(err, services.ErrCinemaHallNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("cinema hall not found by provided ID %v", cinemaHall.CinemaHallID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cinema Hall deleted successfully",
	})
}

func (service *AdminHandler) NewCinemaHallSeatAdmin(c *gin.Context) {
	var newCinemaSeat NewCinemaSeatForm

	if err := c.ShouldBindJSON(&newCinemaSeat); err != nil {
		helpers.RespondWithValidationErrors(c, err, newCinemaSeat)
		return
	}

	err := service.adminCtrl.AddCinemaSeats(newCinemaSeat.SeatRow, newCinemaSeat.SeatNumber, newCinemaSeat.SeatType, newCinemaSeat.HallID)
	if err != nil {
		if errors.Is(err, services.ErrCinemaSeatAlreadyExists) {
			helpers.ClientError(c, http.StatusConflict, fmt.Sprintf("seat already exists with row %s and number %d", newCinemaSeat.SeatRow, newCinemaSeat.SeatNumber))
			return

		}
		if errors.Is(err, services.ErrCinemaHallNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("invalid hall_id: %d. Hall not found", newCinemaSeat.HallID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cinema-Hall Seat added successfully",
	})
}

func (service *AdminHandler) DeleteCinemaHallSeatAdmin(c *gin.Context) {
	var cinemaSeat DeleteCinemaSeatForm

	if err := c.ShouldBindJSON(&cinemaSeat); err != nil {
		helpers.RespondWithValidationErrors(c, err, cinemaSeat)
		return
	}

	err := service.adminCtrl.DeleteCinemaSeat(cinemaSeat.CinemaSeatID)
	if err != nil {
		if errors.Is(err, services.ErrCinemaSeatNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("cinema seat with ID %d not found", cinemaSeat.CinemaSeatID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cinema-Hall Seat deleted successfully",
	})
}

func (service *AdminHandler) AllShowsAdmin(c *gin.Context) {

	allShows, err := service.adminCtrl.FetchAllShowsForAdmin()
	if err != nil {
		if errors.Is(err, services.ErrShowNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "There is no show yet!",
			})
			return
		}
		helpers.ServerError(c, err)
		return
	}

	allMovies, err := service.adminCtrl.FetchAllMovies()
	if err != nil {
		if errors.Is(err, services.ErrAdminPageMovieNotFound) {
			helpers.ClientError(c, http.StatusNotFound, "There is no movie yet!")
			return
		}
		helpers.ServerError(c, err)
		return
	}

	allCinemaHalls, err := service.adminCtrl.FetchAllCinemaHalls()
	if err != nil {
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			helpers.ClientError(c, http.StatusNotFound, "There is no cinema hall yet!")
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"allShows":       allShows,
		"allMovies":      allMovies,
		"allCinemaHalls": allCinemaHalls,
	})
}

func (service *AdminHandler) NewShowAdmin(c *gin.Context) {
	var newShow NewShowForm

	if err := c.ShouldBindJSON(&newShow); err != nil {
		helpers.RespondWithValidationErrors(c, err, newShow)
		return
	}

	err := service.adminCtrl.AddNewShow(newShow.ShowDate, newShow.StartTime, newShow.HallID, newShow.MovieID)
	if err != nil {
		if errors.Is(err, services.ErrCinemaHallNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("cinema hall with ID %d not found", newShow.HallID))
			return
		}
		if errors.Is(err, services.ErrMovieNotFoundByID) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("movie with ID %d not found", newShow.MovieID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "New show added successfully",
	})
}

func (service *AdminHandler) EditShowAdmin(c *gin.Context) {
	var show EditShowForm

	if err := c.ShouldBindJSON(&show); err != nil {
		helpers.RespondWithValidationErrors(c, err, show)
		return
	}

	err := service.adminCtrl.UpdateShow(show.ShowID, show.ShowDate, show.StartTime, show.HallID, show.MovieID)
	if err != nil {
		if errors.Is(err, services.ErrCinemaHallNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("cinema hall with ID %d not found", show.HallID))
			return
		}
		if errors.Is(err, services.ErrMovieNotFoundByID) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("movie with ID %d not found", show.MovieID))
			return
		}
		if errors.Is(err, services.ErrShowNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("show with ID %d not found", show.ShowID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Show updated successfully",
	})
}

func (service *AdminHandler) DeleteShowAdmin(c *gin.Context) {
	var show DeleteShowForm

	if err := c.ShouldBindJSON(&show); err != nil {
		helpers.RespondWithValidationErrors(c, err, show)
		return
	}

	err := service.adminCtrl.DeleteShow(show.ShowID)
	if err != nil {
		if errors.Is(err, services.ErrShowNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("show with ID %d not found", show.ShowID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Show deleted successfully",
	})
}

func (service *AdminHandler) AllShowSeatsAdmin(c *gin.Context) {
	showID, err := helpers.GetParameterFromURL(c, "showID", "invalid show ID provided.")
	if err != nil {
		helpers.ClientError(c, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	allShowSeats, err := service.adminCtrl.FetchAllShowSeats(showID)
	if err != nil {
		if errors.Is(err, services.ErrShowSeatNotFound) {
			helpers.ClientError(c, http.StatusNotFound, "show seats not found")
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"allShowSeats": allShowSeats,
	})
}

func (service *AdminHandler) EditShowSeatPriceAdmin(c *gin.Context) {
	var showSeatPrice EditShowSeatForm

	if err := c.ShouldBindJSON(&showSeatPrice); err != nil {
		helpers.RespondWithValidationErrors(c, err, showSeatPrice)
		return
	}

	err := service.adminCtrl.UpdateShowSeat(showSeatPrice.SeatPrice, showSeatPrice.ShowSeatID)
	if err != nil {
		if errors.Is(err, services.ErrShowSeatNotFound) {
			helpers.ClientError(c, http.StatusNotFound, fmt.Sprintf("show seat with ID %d not found", showSeatPrice.ShowSeatID))
			return
		}
		helpers.ServerError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Seat price updated successfully",
	})
}

// func(service *AdminHandler) AllShowSeatsAdmin(c *gin.Context){}
// func(service *AdminHandler) AllShowSeatsAdmin(c *gin.Context){}
