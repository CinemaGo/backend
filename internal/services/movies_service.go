package services

import (
	"cinemaGo/backend/internal/models"
	"errors"
	"fmt"
)

type MoviesServiceInterface interface {
	FetchAllCaruselImages() ([]models.CarouselImage, error)
	FetchAllMovies() ([]models.AllMovies, error)
	FetchAMovie(movieID int) (models.Movie, error)
	FetchAllActorsCrewsByMovieID(movieID int) ([]models.ActorsCrewsOfMovie, error)
}

type MoviesService struct {
	db models.DBContractMovies
}

func NewMoviesService(db models.DBContractMovies) *MoviesService {
	return &MoviesService{db: db}
}

// FetchAllCaruselImages retrieves all carousel images by calling the
// RetrieveAllCarouselImages method from the database layer.
//
// Parameters:
//   - ms: a pointer to the MoviesService struct, which holds the database connection
//     and is responsible for handling movie-related logic.
//
// Returns:
//   - A slice of models.CarouselImage containing all the carousel images retrieved
//     from the database.
//   - An error if something goes wrong during the fetching process, such as a failure
//     in the database query or scanning the data.
func (ms *MoviesService) FetchAllCaruselImages() ([]models.CarouselImage, error) {
	// Retrieve carousel images using the database method
	carouselImages, err := ms.db.RetrieveAllCarouselImages()
	if err != nil {
		// Return a formatted error if the database retrieval fails
		return nil, fmt.Errorf("error occurred while fetching all carousel images in the service section: %w", err)
	}

	// Return the fetched carousel images
	return carouselImages, nil
}

// FetchAllMovies retrieves all movies from the database and converts the movie ratings
// from a 0-100 scale to a 0-10 scale.
//
// Parameters:
//   - ms: a pointer to the MoviesService struct, which contains the logic for interacting
//     with the database and transforming data.
//
// Returns:
// - A slice of models.AllMovies containing all movie data with the adjusted ratings.
// - An error if there is an issue fetching movies from the database or processing the data.
func (ms *MoviesService) FetchAllMovies() ([]models.AllMovies, error) {
	// Retrieve all movies from the database using the db layer
	movies, err := ms.db.RetrieveAllMovies()
	if err != nil {
		// Return a formatted error if the movie retrieval fails
		return nil, fmt.Errorf("error occured while fetching all movies in the service section: %w", err)
	}

	// Convert movie ratings from a 0-100 scale to a 0-10 scale.
	for i := range movies {
		movies[i].Rating = movies[i].Rating / 10.0
	}

	// Return the slice of movies with adjusted ratings
	return movies, nil
}

// FetchAMovie fetches a movie by its ID from the service layer, converts its rating
// from a 0-100 scale to a 0-10 scale, and handles errors appropriately.
//
// Parameters:
// - id: The unique identifier of the movie to be fetched.
//
// Returns:
// - The Movie struct with movie details if found, including the converted rating.
// - An error if there are issues retrieving the movie or if the movie is not found.
func (ms *MoviesService) FetchAMovie(movieID int) (models.Movie, error) {
	// Call the database method to retrieve the movie by its ID
	movie, err := ms.db.RetrieveAMovie(movieID)
	if err != nil {
		// If the movie is not found, return a specific error
		if errors.Is(err, models.ErrMovieNotFoundByID) {
			return models.Movie{}, ErrMovieNotFoundByID
		}

		// For other errors, wrap and return a descriptive error
		return models.Movie{}, fmt.Errorf("error occurred while fetching a movie by id in the service section: %w", err)
	}

	// Convert the movie rating from a 0-100 scale to a 0-10 scale
	movie.Rating = float32(movie.Rating) / 10.0

	// Return the movie with the adjusted rating
	return movie, nil
}

// FetchAllActorsCrewsByMovieID fetches all actors and crew members associated with a movie
// based on the provided movie ID. It acts as a service layer method that interfaces with
// the database layer to retrieve the data.
//
// Parameters:
// - movieID: The ID of the movie for which actors and crew need to be fetched.
//
// Returns:
// - A slice of ActorsCrewsOfMovie representing the details of the actors/crew associated with the movie.
// - An error if there are issues with fetching the data.
func (ms *MoviesService) FetchAllActorsCrewsByMovieID(movieID int) ([]models.ActorsCrewsOfMovie, error) {
	// Call the database function to retrieve actors and crew for the given movie ID
	allActorsCrew, err := ms.db.RetrieveAllActorsCrewsByMovieID(movieID)
	if err != nil {
		// If there is an error, wrap it with context and return
		return nil, fmt.Errorf("error occurred while fetching all actors and crews in the service section: %w", err)
	}

	// Return the list of actors/crew if no errors occurred
	return allActorsCrew, nil
}
