package services

import (
	"cinemaGo/backend/internal/models"
	"fmt"
)

type MoviesServiceInterface interface {
	FetchAllCaruselImages() ([]models.CarouselImage, error)
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
