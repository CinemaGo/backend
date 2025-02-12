package services

import (
	"cinemaGo/backend/internal/models"
	"errors"
	"fmt"
	"time"
)

type AdminServiceInterface interface {
	CreateNewCarouselImage(imageURL, title, description string, orderPriority int) error
	FetchAllCarouselImages() ([]models.CarouselImageForAdmin, error)
	UpdateCarouselImages(imageURL, title, description string, orderPriority int, carouselImageID int) error
	DeleteCarouselImages(carouselImageID int) error

	AddNewMovie(title string, description, genre, language, trailerURL, posterURL string, rating float32, ratingProvider string, duration int, releaseDate, ageLimit string) error
	FetchAllMovies() ([]models.AllMoviesForAdmin, error)
	FetchAMovie(movieID int) (models.MovieForAdmin, error)
	UpdateMovieInfo(movieID int, title, description, genre, language, trailerURL, posterURL string, rating float32, ratingProvider string, duration int, relaseDate string, ageLimit string) error
	DeleteMovie(movieID int) error

	AddActorsCrew(fullName, imageURL, occupation, roleDescription string, bornDate time.Time, birthplace, about string, isActor bool, movieID int) error
	FetchAllActorsCrewByMovieID(movieID int) ([]models.AllActorCrewForAdmin, error)
	FetchAnActorCrew(actorCrewID int) (models.ActorCrewForAdmin, error)
	UpdateActorCrewInfo(fullName, imageURL, occupation, roleDescription string, bornDate time.Time, birthplace, about string, isActor bool, actorCrewID int) error
	DeleteActorCrew(actorCrewID int) error

	AddNewCinemaHall(hallName, hallType string, capacity int) error
	FetchAllCinemaHalls() ([]models.CinemaHallForAdmin, error)
	FetchCinemaHallInfo(cinemaHallID int) (models.CinemaHallForAdmin, error)
	UpdateCinemaHallInfo(hallName, hallType string, capacity, cinemaHallID int) error
	DeleteCinemaHall(cinemaHallID int) error

	AddCinemaSeats(seatRow string, seatNumber int, seatType string, hallID int) error
	FetchALLCinemaSeatsByHallID(hallID int) ([]models.CinemaSeatForAdmin, error)
	DeleteCinemaSeat(cinemaSeatID int) error

	AddNewShow(showDate, startTime time.Time, hallID int, movieID int) error
	FetchAllShowsForAdmin() ([]models.ShowForAdmin, error)
	UpdateShow(showID int, showDate, startTime time.Time, hallID int, movieID int) error
	DeleteShow(showID int) error

	FetchAllShowSeats(showID int) ([]models.ShowSeatForAdmin, error)
	UpdateShowSeat(seatPrice float32, showSeatID int) error
}

type AdminService struct {
	db models.DBContractAdminCtrl
}

func NewAdminService(db models.DBContractAdminCtrl) *AdminService {
	return &AdminService{db: db}
}

func (as *AdminService) CreateNewCarouselImage(imageURL, title, description string, orderPriority int) error {
	err := as.db.InsertNewCarouselImages(imageURL, title, description, orderPriority)
	if err != nil {
		return fmt.Errorf("error occured while creating new carousel images in the service section: %w", err)
	}

	return nil
}

func (as *AdminService) FetchAllCarouselImages() ([]models.CarouselImageForAdmin, error) {

	carouselImages, err := as.db.RetrieveCarouselImagesForAdmin()
	if err != nil {
		if errors.Is(err, models.ErrAdminPageCarouselImagesNotFound) {
			return nil, ErrAdminPageCarouselImagesNotFound
		}
		return nil, fmt.Errorf("error occured while fetfetching carousel images in the service serction: %w", err)
	}

	return carouselImages, nil

}

func (as *AdminService) UpdateCarouselImages(imageURL, title, description string, orderPriority int, carouselImageID int) error {
	err := as.db.UpdateCarouselImagesByID(imageURL, title, description, orderPriority, carouselImageID)
	if err != nil {
		if errors.Is(err, models.ErrActorCrewNotFoundByID) {
			return ErrAdminPageCarouselImagesNotFound
		}
		return fmt.Errorf("error occured while updating carousel image data: %w", err)
	}
	return nil
}

func (as *AdminService) DeleteCarouselImages(carouselImageID int) error {
	err := as.db.DeleteCarouselImagesByID(carouselImageID)
	if err != nil {
		if errors.Is(err, models.ErrActorCrewNotFoundByID) {
			return ErrAdminPageCarouselImagesNotFound
		}
		return fmt.Errorf("error occured while deleting carousel image data: %w", err)
	}

	return nil
}

func (as *AdminService) AddNewMovie(title string, description, genre, language, trailerURL, posterURL string, rating float32, ratingProvider string, duration int, releaseDate, ageLimit string) error {

	rating = rating * 10

	err := as.db.InsertNewMovie(title, description, genre, language, trailerURL, posterURL, int(rating), ratingProvider, duration, releaseDate, ageLimit)
	if err != nil {
		return fmt.Errorf("error occured while adding new movie in the service section: %w", err)
	}

	return nil
}

func (as *AdminService) FetchAllMovies() ([]models.AllMoviesForAdmin, error) {
	allMovies, err := as.db.RetrieveAllMoviesForAdmin()
	if err != nil {
		if errors.Is(err, models.ErrAdminPageMovieNotFound) {
			return nil, ErrAdminPageMovieNotFound
		}
		return nil, fmt.Errorf("error occured while fetching all movies for admin: %w", err)
	}

	return allMovies, nil
}

func (as *AdminService) FetchAMovie(movieID int) (models.MovieForAdmin, error) {
	movie, err := as.db.RetrieveAMovieForAdmin(movieID)
	if err != nil {
		if errors.Is(err, models.ErrAdminPageMovieNotFound) {
			return models.MovieForAdmin{}, ErrAdminPageMovieNotFound
		}
		return models.MovieForAdmin{}, fmt.Errorf("error occured while fetching movie for admin page: %w", err)
	}
	movie.Rating = float32(movie.Rating) / 10.0

	return movie, nil
}

func (as *AdminService) UpdateMovieInfo(movieID int, title, description, genre, language, trailerURL, posterURL string, rating float32, ratingProvider string, duration int, relaseDate string, ageLimit string) error {

	rating = rating * 10

	err := as.db.UpdateMovieInfoForAdminByMovieID(movieID, title, description, genre, language, trailerURL, posterURL, rating, ratingProvider, duration, relaseDate, ageLimit)
	if err != nil {
		if errors.Is(err, models.ErrAdminPageMovieNotFound) {
			return ErrAdminPageMovieNotFound
		}
		return fmt.Errorf("error occured while updating movie data: %w", err)
	}
	return nil
}

func (as *AdminService) DeleteMovie(movieID int) error {
	err := as.db.DeleteMovieByMovieID(movieID)
	if err != nil {
		if errors.Is(err, models.ErrAdminPageMovieNotFound) {
			return ErrAdminPageMovieNotFound
		}
		return fmt.Errorf("error occured while deleting movie: %w", err)
	}

	return nil
}

func (as *AdminService) AddActorsCrew(fullName, imageURL, occupation, roleDescription string, bornDate time.Time, birthplace, about string, isActor bool, movieID int) error {

	formattedBornDate := bornDate.Format("2006-01-02")

	actorCrewID, err := as.db.InsertActorsCrew(fullName, imageURL, occupation, roleDescription, formattedBornDate, birthplace, about, isActor)
	if err != nil {
		return fmt.Errorf("error occured while adding new actorCrew: %w", err)
	}

	err = as.db.InsertMovieActorCrew(movieID, actorCrewID)
	if err != nil {
		return fmt.Errorf("error occured while adding movie actorCrew: %w", err)
	}

	return nil
}

func (as *AdminService) FetchAllActorsCrewByMovieID(movieID int) ([]models.AllActorCrewForAdmin, error) {

	allActorCrew, err := as.db.RetrieveAllActorsCrewByMovieID(movieID)
	if err != nil {
		if errors.Is(err, models.ErrActorCrewNotFound) {
			return nil, ErrActorCrewNotFound
		}
		return nil, fmt.Errorf("error occured while fetching all actorsCrew: %w", err)
	}

	return allActorCrew, nil
}

func (as *AdminService) FetchAnActorCrew(actorCrewID int) (models.ActorCrewForAdmin, error) {
	actorCrew, err := as.db.RetrieveActorCrew(actorCrewID)
	if err != nil {
		if errors.Is(err, models.ErrActorCrewNotFound) {
			return models.ActorCrewForAdmin{}, ErrActorCrewNotFound
		}
		return models.ActorCrewForAdmin{}, fmt.Errorf("error occured while fetching actor/crew info: %w", err)
	}

	return actorCrew, nil
}

func (as *AdminService) UpdateActorCrewInfo(fullName, imageURL, occupation, roleDescription string, bornDate time.Time, birthplace, about string, isActor bool, actorCrewID int) error {

	formattedBornDate := bornDate.Format("2006-01-02")

	err := as.db.UpdateActorCrewInformation(fullName, imageURL, occupation, roleDescription, formattedBornDate, birthplace, about, isActor, actorCrewID)
	if err != nil {
		if errors.Is(err, models.ErrActorCrewNotFound) {
			return ErrActorCrewNotFound
		}
		return fmt.Errorf("error occured while updating actor/crew info: %w", err)
	}

	return nil
}

func (as *AdminService) DeleteActorCrew(actorCrewID int) error {
	err := as.db.DeleteActorCrewByID(actorCrewID)
	if err != nil {
		if errors.Is(err, models.ErrActorCrewNotFound) {
			return ErrActorCrewNotFound
		}
		return fmt.Errorf("error occured while deleting actor/crew info: %w", err)
	}

	return nil
}

func (as *AdminService) AddNewCinemaHall(hallName, hallType string, capacity int) error {
	err := as.db.InsertNewCinemaHall(hallName, hallType, capacity)
	if err != nil {
		return fmt.Errorf("error occured while adding new hall: %w", err)
	}

	return nil
}

func (as *AdminService) FetchAllCinemaHalls() ([]models.CinemaHallForAdmin, error) {
	allCinemaHalls, err := as.db.RetrieveAllCinemaHallsForAdmin()
	if err != nil {
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return nil, ErrCinemaHallNotFound
		}
		return nil, fmt.Errorf("error occured while fetching all cinema halls: %w", err)
	}

	return allCinemaHalls, nil
}

func (as *AdminService) FetchCinemaHallInfo(cinemaHallID int) (models.CinemaHallForAdmin, error) {
	cinemaHall, err := as.db.RetrieveCinemaHallInfoByID(cinemaHallID)
	if err != nil {
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return models.CinemaHallForAdmin{}, ErrCinemaHallNotFound
		}
		return models.CinemaHallForAdmin{}, fmt.Errorf("error occured while fetching cinema hall: %w", err)
	}

	return cinemaHall, nil
}

func (as *AdminService) UpdateCinemaHallInfo(hallName, hallType string, capacity, cinemaHallID int) error {
	err := as.db.UpdateCinemaHallInfoByID(hallName, hallType, capacity, cinemaHallID)
	if err != nil {
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return ErrCinemaHallNotFound
		}
		return fmt.Errorf("error occured while updating cinema hall data: %w", err)
	}

	return nil
}

func (as *AdminService) DeleteCinemaHall(cinemaHallID int) error {
	err := as.db.DeleteCinemaHallByID(cinemaHallID)
	if err != nil {
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return ErrCinemaHallNotFound
		}
		return fmt.Errorf("error occured while deleting cinema hall: %w", err)
	}

	return nil
}

func (as *AdminService) AddCinemaSeats(seatRow string, seatNumber int, seatType string, hallID int) error {
	err := as.db.InsertCinemaSeats(seatRow, seatNumber, seatType, hallID)
	if err != nil {
		if errors.Is(err, models.ErrCinemaSeatAlreadyExists) {
			return ErrCinemaSeatAlreadyExists
		}

		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return ErrCinemaHallNotFound
		}

		return fmt.Errorf("error occured while adding new cinema seat: %w", err)
	}

	return nil
}

func (as *AdminService) FetchALLCinemaSeatsByHallID(hallID int) ([]models.CinemaSeatForAdmin, error) {

	allCinemaSeat, err := as.db.RetrieveALLCinemaSeatsByHallID(hallID)
	if err != nil {
		if errors.Is(err, models.ErrCinemaSeatNotFound) {
			return nil, ErrCinemaSeatNotFound
		}
		return nil, fmt.Errorf("error occured while fetching all cinema seat by hall id: %w", err)
	}

	return allCinemaSeat, nil
}

func (as *AdminService) DeleteCinemaSeat(cinemaSeatID int) error {
	err := as.db.DeleteCinemaSeatByID(cinemaSeatID)
	if err != nil {
		if errors.Is(err, models.ErrCinemaSeatNotFound) {
			return ErrCinemaSeatNotFound
		}
		return fmt.Errorf("error occured while deleting cinema seat: %w", err)
	}

	return nil
}

func (as *AdminService) AddNewShow(showDate, startTime time.Time, hallID int, movieID int) error {

	formattedDate := showDate.Format("2006-01-02")
	formattedTime := startTime.Format("15:04:05")

	showID, err := as.db.InsertNewShow(formattedDate, formattedTime, hallID, movieID)
	if err != nil {
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return ErrCinemaHallNotFound
		}
		if errors.Is(err, models.ErrMovieNotFoundByID) {
			return ErrMovieNotFoundByID
		}

		return fmt.Errorf("error occured while addning new show: %w", err)
	}

	allCinemaSeats, err := as.db.RetrieveALLCinemaSeatsByHallID(hallID)
	if err != nil {
		if errors.Is(err, models.ErrCinemaSeatNotFound) {
			return ErrCinemaSeatNotFound
		}
		return fmt.Errorf("error occured while retrieving cinema hall seats")
	}

	for _, cinemaSeat := range allCinemaSeats {
		err = as.db.InsertNewShowSeat("Available", 0, cinemaSeat.CinemaSeatID, showID)
		if err != nil {
			if errors.Is(err, models.ErrShowSeatNotFound) {
				return ErrShowSeatNotFound
			}
			return fmt.Errorf("error occured while inserting new show seats: %w", err)
		}
	}

	return nil
}

func (as *AdminService) FetchAllShowsForAdmin() ([]models.ShowForAdmin, error) {
	allShows, err := as.db.RetrieveAllShowsForAdmin()
	if err != nil {
		if errors.Is(err, models.ErrShowNotFound) {
			return nil, ErrShowNotFound
		}
		return nil, fmt.Errorf("error occured while fetching all shows: %w", err)
	}

	return allShows, nil
}

func (as *AdminService) UpdateShow(showID int, showDate, startTime time.Time, hallID int, movieID int) error {

	formattedDate := showDate.Format("2006-01-02")
	formattedTime := startTime.Format("15:04:05")

	err := as.db.UpdateShowByID(showID, formattedDate, formattedTime, hallID, movieID)
	if err != nil {
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return ErrCinemaHallNotFound
		}
		if errors.Is(err, models.ErrMovieNotFoundByID) {
			return ErrMovieNotFoundByID
		}
		if errors.Is(err, models.ErrShowNotFound) {
			return ErrShowNotFound
		}
		return fmt.Errorf("error occured while updating show data: %w", err)
	}

	return nil
}

func (as *AdminService) DeleteShow(showID int) error {
	err := as.db.DeleteShowByID(showID)
	if err != nil {
		if errors.Is(err, models.ErrShowNotFound) {
			return ErrShowNotFound
		}
		return fmt.Errorf("error occured while deleting show: %w", err)
	}

	return nil
}

func (as *AdminService) FetchAllShowSeats(showID int) ([]models.ShowSeatForAdmin, error) {

	allShowSeats, err := as.db.RetrieveAllShowSeats(showID)
	if err != nil {
		if errors.Is(err, models.ErrShowSeatNotFound) {
			return nil, ErrShowSeatNotFound
		}
		return nil, fmt.Errorf("error occured while fetching all show seats: %w", err)
	}

	return allShowSeats, nil
}

func (as *AdminService) UpdateShowSeat(seatPrice float32, showSeatID int) error {

	seatPrice = seatPrice * 10

	err := as.db.UpdateShowSeatByID(int(seatPrice), showSeatID)
	if err != nil {
		if errors.Is(err, models.ErrShowSeatNotFound) {
			return ErrShowSeatNotFound
		}
		return fmt.Errorf("error occured while updating seat price: %w", err)
	}

	return nil
}
