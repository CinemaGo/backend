package services

import (
	"cinemaGo/backend/internal/models"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type AdminServiceInterface interface {
	CreateNewCarouselImage(imageURL, title, description string, orderPriority int) error
	FetchAllCarouselImages() ([]models.CarouselImageForAdmin, error)
	UpdateCarouselImages(imageURL, title, description string, orderPriority string, carouselImageID int) error
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

// CreateNewCarouselImage creates a new carousel image by calling the InsertNewCarouselImages method
// from the database layer.
// Parameters:
//   - imageURL (string): URL of the image to be added to the carousel
//   - title (string): Title for the carousel image
//   - description (string): Description for the carousel image
//   - orderPriority (int): Order priority for the carousel image (used for sorting)
//
// Returns:
//   - error: Returns an error if something goes wrong during the insert operation
func (as *AdminService) CreateNewCarouselImage(imageURL, title, description string, orderPriority int) error {
	// Call the database method to insert the carousel image
	err := as.db.InsertNewCarouselImages(imageURL, title, description, orderPriority)
	if err != nil {
		// If there is an error, return it wrapped in a descriptive message
		return fmt.Errorf("error occurred while creating new carousel images in the service section: %w", err)
	}

	// If everything goes well, return nil (indicating success)
	return nil
}

// FetchAllCarouselImages retrieves all carousel images for the admin page
// from the database and handles errors appropriately.
// Returns:
//   - []models.CarouselImageForAdmin: A slice of carousel images for admin.
//   - error: Returns an error if anything goes wrong during the fetch operation.
func (as *AdminService) FetchAllCarouselImages() ([]models.CarouselImageForAdmin, error) {

	// Call the database method to retrieve all carousel images
	carouselImages, err := as.db.RetrieveCarouselImagesForAdmin()
	if err != nil {
		// If the error is that no carousel images are found, return the custom error
		if errors.Is(err, models.ErrAdminPageCarouselImagesNotFound) {
			return nil, ErrAdminPageCarouselImagesNotFound
		}
		// Otherwise, wrap the error with additional context and return it
		return nil, fmt.Errorf("error occurred while fetching carousel images in the service section: %w", err)
	}

	// Return the retrieved carousel images if no errors occurred
	return carouselImages, nil
}

// UpdateCarouselImages updates an existing carousel image in the admin panel.
//
// Parameters:
//   - imageURL (string): The URL of the new image to be used for the carousel.
//   - title (string): The new title to be associated with the carousel image.
//   - description (string): The new description for the carousel image.
//   - orderPriority (int): The new display order priority for the carousel image.
//   - carouselImageID (int): The unique identifier of the carousel image to be updated.
//
// Returns:
//   - error: Returns nil if the update is successful, otherwise returns an error
//     with a detailed message of what went wrong.
func (as *AdminService) UpdateCarouselImages(imageURL, title, description string, orderPriority string, carouselImageID int) error {
	intOrderPriority, err := strconv.Atoi(orderPriority)
	if err != nil {
		return fmt.Errorf("error occurred while converting orderPriority from string to int: %w", err)
	}

	// Attempt to update the carousel image by calling the database method.
	err = as.db.UpdateCarouselImagesByID(imageURL, title, description, intOrderPriority, carouselImageID)
	if err != nil {
		// If the error indicates that the carousel image was not found, return a custom error.
		if errors.Is(err, models.ErrActorCrewNotFoundByID) {
			return ErrAdminPageCarouselImagesNotFound
		}
		// If the update fails, return a wrapped error with details.
		return fmt.Errorf("error occurred while updating carousel image data: %w", err)
	}
	// Return nil if the update was successful.
	return nil
}

// DeleteCarouselImages deletes a carousel image by its unique ID in the admin panel.
//
// Parameters:
//   - carouselImageID (int): The unique identifier of the carousel image to be deleted.
//
// Returns:
//   - error: Returns nil if the deletion is successful, otherwise returns an error
//     with a detailed message of what went wrong.
func (as *AdminService) DeleteCarouselImages(carouselImageID int) error {
	// Attempt to delete the carousel image by calling the database method.
	err := as.db.DeleteCarouselImagesByID(carouselImageID)
	if err != nil {
		// If the error indicates that the carousel image was not found, return a custom error.
		if errors.Is(err, models.ErrAdminPageCarouselImagesNotFound) {
			return ErrAdminPageCarouselImagesNotFound
		}
		// If the deletion fails, return a wrapped error with additional context.
		return fmt.Errorf("error occurred while deleting carousel image data: %w", err)
	}

	// Return nil if the deletion was successful.
	return nil
}

// AddNewMovie adds a new movie to the database by calling the appropriate method in the database layer.
//
// Parameters:
//   - title (string): The title of the movie.
//   - description (string): A brief description or plot summary of the movie.
//   - genre (string): The genre or category the movie belongs to (e.g., Action, Drama).
//   - language (string): The language in which the movie is made.
//   - trailerURL (string): The URL link to the movie's trailer.
//   - posterURL (string): The URL link to the movie's poster image.
//   - rating (float32): The movie's rating (scaled from 0 to 10, later multiplied by 10 for storage).
//   - ratingProvider (string): The name or source providing the movie's rating (e.g., IMDB).
//   - duration (int): The duration of the movie in minutes.
//   - releaseDate (string): The release date of the movie in string format (e.g., "YYYY-MM-DD").
//   - ageLimit (string): The age restriction or rating for the movie (e.g., "PG-13").
//
// Returns:
//   - error: Returns nil if the movie was successfully added, or an error if something went wrong.
//     The error provides additional context, including details about where the process failed.
func (as *AdminService) AddNewMovie(title string, description, genre, language, trailerURL, posterURL string, rating float32, ratingProvider string, duration int, releaseDate, ageLimit string) error {
	// Multiply rating by 10 to scale it from 0-10 to 0-100 (as per system requirements).
	rating = rating * 10

	// Attempt to insert the new movie into the database.
	err := as.db.InsertNewMovie(title, description, genre, language, trailerURL, posterURL, int(rating), ratingProvider, duration, releaseDate, ageLimit)
	if err != nil {
		// If an error occurs, return it with additional context.
		return fmt.Errorf("error occurred while adding new movie in the service section: %w", err)
	}

	// Return nil if the movie was successfully added.
	return nil
}

// FetchAllMovies retrieves all movies from the database for the admin page.
//
// This function calls the database layer to fetch a list of all movies available in the system,
// specifically designed for the admin interface.
//
// Parameters:
//   - None
//
// Returns:
//   - ([]models.AllMoviesForAdmin): A slice containing all movie data for the admin interface.
//   - error: Returns nil if the movies are successfully retrieved. If an error occurs, it will
//     provide more context about the failure (e.g., if no movies are found, or if there's a database issue).
func (as *AdminService) FetchAllMovies() ([]models.AllMoviesForAdmin, error) {
	// Attempt to retrieve all movies for the admin page from the database.
	allMovies, err := as.db.RetrieveAllMoviesForAdmin()
	if err != nil {
		// If no movies are found, return a custom error indicating no movies exist.
		if errors.Is(err, models.ErrAdminPageMovieNotFound) {
			return nil, ErrAdminPageMovieNotFound
		}
		// Wrap the error with more context if something goes wrong.
		return nil, fmt.Errorf("error occurred while fetching all movies for admin: %w", err)
	}

	// Return the list of movies if no error occurs.
	return allMovies, nil
}

// FetchAMovie retrieves a specific movie's details from the database for the admin page.
//
// This function queries the database to fetch detailed information about a movie based on its ID
// and adjusts the movie's rating before returning it.
//
// Parameters:
//   - movieID (int): The unique identifier for the movie to retrieve.
//
// Returns:
//   - models.MovieForAdmin: The details of the movie for the admin interface, including title,
//     description, genre, language, rating, etc.
//   - error: Returns nil if the movie is successfully retrieved. If an error occurs (e.g., movie not found),
//     a descriptive error is returned.
func (as *AdminService) FetchAMovie(movieID int) (models.MovieForAdmin, error) {
	// Attempt to retrieve a movie's details by its ID for the admin page from the database.
	movie, err := as.db.RetrieveAMovieForAdmin(movieID)
	if err != nil {
		// If the movie is not found, return a custom error indicating the movie doesn't exist.
		if errors.Is(err, models.ErrAdminPageMovieNotFound) {
			return models.MovieForAdmin{}, ErrAdminPageMovieNotFound
		}
		// Wrap any other errors with more context and return them.
		return models.MovieForAdmin{}, fmt.Errorf("error occurred while fetching movie for admin page: %w", err)
	}

	// Adjust the rating to a decimal format (scaled down by a factor of 10).
	movie.Rating = float32(movie.Rating) / 10.0

	// Return the movie details along with a successfully fetched status.
	return movie, nil
}

// UpdateMovieInfo updates the details of a movie in the database for the admin page.
//
// This function allows the admin to update movie details like title, description, genre, language,
// trailer URL, poster URL, rating, and more. The rating is scaled from a decimal to an integer representation
// before updating in the database.
//
// Parameters:
//   - movieID (int): The unique identifier of the movie to be updated.
//   - title (string): The title of the movie.
//   - description (string): A detailed description of the movie.
//   - genre (string): The genre of the movie (e.g., Drama, Comedy, etc.).
//   - language (string): The language of the movie.
//   - trailerURL (string): A URL to the movie's trailer.
//   - posterURL (string): A URL to the movie's poster image.
//   - rating (float32): The rating of the movie (out of 10), which is multiplied by 10 before storing in the database.
//   - ratingProvider (string): The source of the rating (e.g., Rotten Tomatoes, IMDb).
//   - duration (int): The duration of the movie in minutes.
//   - releaseDate (string): The release date of the movie in YYYY-MM-DD format.
//   - ageLimit (string): The age rating of the movie (e.g., PG, R, etc.).
//
// Returns:
//   - error: Returns nil if the update is successful. If an error occurs (such as the movie not being found),
//     it will return a descriptive error.
func (as *AdminService) UpdateMovieInfo(movieID int, title, description, genre, language, trailerURL, posterURL string, rating float32, ratingProvider string, duration int, releaseDate string, ageLimit string) error {
	// Scale the rating by 10 to store it in an integer format (usually 1-100 in the database).
	rating = rating * 10

	// Attempt to update the movie details in the database.
	err := as.db.UpdateMovieInfoForAdminByMovieID(movieID, title, description, genre, language, trailerURL, posterURL, rating, ratingProvider, duration, releaseDate, ageLimit)
	if err != nil {
		// If the movie is not found in the database, return a specific error indicating the movie doesn't exist.
		if errors.Is(err, models.ErrAdminPageMovieNotFound) {
			return ErrAdminPageMovieNotFound
		}
		// Return any other errors, wrapped with additional context to explain where the error occurred.
		return fmt.Errorf("error occurred while updating movie data: %w", err)
	}

	// If the update was successful, return nil to indicate no error.
	return nil
}

// DeleteMovie removes a movie from the database based on its ID.
//
// This function allows the admin to delete a movie by providing its unique identifier (movieID).
// If the movie does not exist, a specific error is returned. Any other errors during deletion
// are wrapped with a message providing context for easier debugging.
//
// Parameters:
//   - movieID (int): The unique identifier of the movie to be deleted.
//
// Returns:
//   - error: Returns nil if the deletion is successful. If an error occurs (such as the movie not being found),
//     it will return a descriptive error indicating the problem.
func (as *AdminService) DeleteMovie(movieID int) error {
	// Attempt to delete the movie from the database using its unique ID.
	err := as.db.DeleteMovieByMovieID(movieID)
	if err != nil {
		// If the movie is not found in the database, return a specific error indicating that the movie doesn't exist.
		if errors.Is(err, models.ErrAdminPageMovieNotFound) {
			return ErrAdminPageMovieNotFound
		}
		// Return any other errors, wrapping them with additional context to explain where the error occurred.
		return fmt.Errorf("error occurred while deleting movie: %w", err)
	}

	// If the deletion was successful, return nil to indicate no error.
	return nil
}

// AddActorsCrew adds a new actor/crew member and associates them with a movie.
//
// This function takes the details of an actor or crew member and their associated movie,
// then inserts the new actor/crew member into the database. It also creates a relationship
// between the actor/crew member and the movie. If any of the database operations fail,
// a descriptive error is returned.
//
// Parameters:
//   - fullName (string): The full name of the actor or crew member.
//   - imageURL (string): A URL to the actor's or crew member's image.
//   - occupation (string): The occupation or role of the individual (e.g., Actor, Director).
//   - roleDescription (string): A description of the role or work of the actor/crew member in the movie.
//   - bornDate (time.Time): The birthdate of the actor/crew member.
//   - birthplace (string): The birthplace of the actor/crew member.
//   - about (string): Additional details or biography about the actor/crew member.
//   - isActor (bool): Flag indicating whether the individual is an actor (true) or a crew member (false).
//   - movieID (int): The ID of the movie with which the actor/crew member will be associated.
//
// Returns:
//   - error: Returns nil if the operation is successful. If any error occurs (e.g., inserting actor/crew member
//     or associating them with the movie), an error with context is returned.
func (as *AdminService) AddActorsCrew(fullName, imageURL, occupation, roleDescription string, bornDate time.Time, birthplace, about string, isActor bool, movieID int) error {

	// Format the birthdate into a string for database insertion.
	formattedBornDate := bornDate.Format("2006-01-02")

	// Insert the new actor/crew member into the database.
	actorCrewID, err := as.db.InsertActorsCrew(fullName, imageURL, occupation, roleDescription, formattedBornDate, birthplace, about, isActor)
	if err != nil {
		// Return a descriptive error if the insertion of the actor/crew member fails.
		return fmt.Errorf("error occurred while adding new actor/crew member: %w", err)
	}

	// Associate the newly added actor/crew member with the movie.
	err = as.db.InsertMovieActorCrew(movieID, actorCrewID)
	if err != nil {
		// Return a descriptive error if the association of the actor/crew member with the movie fails.
		return fmt.Errorf("error occurred while associating actor/crew member with movie: %w", err)
	}

	// If both operations are successful, return nil indicating no error.
	return nil
}

// FetchAllActorsCrewByMovieID retrieves all actors and crew members associated with a specific movie.
//
// This function fetches the list of all actors and crew members associated with the provided movie ID.
// If there are no actors or crew members found for the movie, it will return an error. If there are any
// issues while fetching the data from the database, a descriptive error is returned.
//
// Parameters:
//   - movieID (int): The ID of the movie for which the actors and crew need to be fetched.
//
// Returns:
//   - ([]models.AllActorCrewForAdmin, error): A slice of `AllActorCrewForAdmin` containing the details
//     of the actors and crew associated with the movie. If an error occurs, an error is returned with
//     context about the failure.
func (as *AdminService) FetchAllActorsCrewByMovieID(movieID int) ([]models.AllActorCrewForAdmin, error) {

	// Retrieve all actors and crew for the specified movieID.
	allActorCrew, err := as.db.RetrieveAllActorsCrewByMovieID(movieID)
	if err != nil {
		// If actors/crew are not found, return a specific error.
		if errors.Is(err, models.ErrActorCrewNotFound) {
			return nil, ErrActorCrewNotFound
		}
		// Return a descriptive error if any other error occurs while fetching the data.
		return nil, fmt.Errorf("error occurred while fetching all actors and crew: %w", err)
	}

	// Return the list of actors and crew.
	return allActorCrew, nil
}

// FetchAnActorCrew retrieves the details of a specific actor or crew member by their ID.
//
// This function queries the database to fetch the details of an actor or crew member associated
// with the provided actorCrewID. If the actor or crew member is not found, a specific error is returned.
// In case of any other database-related issues, a wrapped error with context is returned.
//
// Parameters:
//   - actorCrewID (int): The ID of the actor or crew member to fetch details for.
//
// Returns:
//   - (models.ActorCrewForAdmin, error): The `ActorCrewForAdmin` struct containing details about the actor or crew.
//     If an error occurs, an error is returned with context to help identify the cause of the failure.
func (as *AdminService) FetchAnActorCrew(actorCrewID int) (models.ActorCrewForAdmin, error) {
	// Retrieve the actor or crew member details from the database using the provided actorCrewID.
	actorCrew, err := as.db.RetrieveActorCrew(actorCrewID)
	if err != nil {
		// If the actor or crew member is not found, return a specific error.
		if errors.Is(err, models.ErrActorCrewNotFound) {
			return models.ActorCrewForAdmin{}, ErrActorCrewNotFound
		}
		// Return a wrapped error with context if any other error occurs while fetching the actor/crew data.
		return models.ActorCrewForAdmin{}, fmt.Errorf("error occurred while fetching actor/crew info: %w", err)
	}

	// Return the retrieved actor or crew member data.
	return actorCrew, nil
}

// UpdateActorCrewInfo updates the information of an actor or crew member in the database.
//
// This function allows for updating various details (such as name, occupation, etc.) of an actor or crew
// member using the provided actorCrewID. If the actor or crew member is not found, a specific error is returned.
// If any error occurs during the update process, a wrapped error with context is returned.
//
// Parameters:
//   - fullName (string): The full name of the actor or crew member.
//   - imageURL (string): The URL of the actor's or crew member's image.
//   - occupation (string): The occupation of the actor or crew member (e.g., actor, director, etc.).
//   - roleDescription (string): A description of the role the actor or crew member played.
//   - bornDate (time.Time): The birthdate of the actor or crew member.
//   - birthplace (string): The birthplace of the actor or crew member.
//   - about (string): A short biography or description about the actor or crew member.
//   - isActor (bool): A boolean indicating whether the actorCrewID is for an actor (true) or crew member (false).
//   - actorCrewID (int): The unique ID of the actor or crew member to update.
//
// Returns:
//   - error: Returns `nil` if the update is successful. Otherwise, it returns an error explaining what went wrong.
func (as *AdminService) UpdateActorCrewInfo(fullName, imageURL, occupation, roleDescription string, bornDate time.Time, birthplace, about string, isActor bool, actorCrewID int) error {
	// Format the birthdate to match the database's expected date format (YYYY-MM-DD).
	formattedBornDate := bornDate.Format("2006-01-02")

	// Attempt to update the actor/crew information in the database.
	err := as.db.UpdateActorCrewInformation(fullName, imageURL, occupation, roleDescription, formattedBornDate, birthplace, about, isActor, actorCrewID)
	if err != nil {
		// If the actor or crew member is not found in the database, return a specific error.
		if errors.Is(err, models.ErrActorCrewNotFound) {
			return ErrActorCrewNotFound
		}
		// Return a wrapped error with context if any other error occurs during the update process.
		return fmt.Errorf("error occurred while updating actor/crew info: %w", err)
	}

	// Return nil if the update was successful.
	return nil
}

// DeleteActorCrew deletes an actor or crew member from the database using their unique ID.
//
// This function attempts to remove an actor or crew member's information from the database
// using the provided actorCrewID. If the actor or crew member is not found, a specific error is returned.
// If any other error occurs during the deletion process, a wrapped error is returned.
//
// Parameters:
//   - actorCrewID (int): The unique ID of the actor or crew member to delete.
//
// Returns:
//   - error: Returns `nil` if the deletion is successful. Otherwise, it returns an error explaining what went wrong.
func (as *AdminService) DeleteActorCrew(actorCrewID int) error {
	// Attempt to delete the actor/crew member from the database using the provided actorCrewID.
	err := as.db.DeleteActorCrewByID(actorCrewID)
	if err != nil {
		// If the actor/crew member is not found, return a specific error.
		if errors.Is(err, models.ErrActorCrewNotFound) {
			return ErrActorCrewNotFound
		}
		// Return a wrapped error with context if any other error occurs during the deletion process.
		return fmt.Errorf("error occurred while deleting actor/crew info: %w", err)
	}

	// Return nil if the deletion was successful.
	return nil
}

// AddNewCinemaHall adds a new cinema hall to the database.
//
// This function attempts to insert a new cinema hall into the database with the provided details.
// If an error occurs during the insertion, it returns a wrapped error. If successful, it returns nil.
//
// Parameters:
//   - hallName (string): The name of the new cinema hall.
//   - hallType (string): The type of the cinema hall (e.g., IMAX, 3D, Regular, etc.).
//   - capacity (int): The seating capacity of the new cinema hall.
//
// Returns:
//   - error: Returns `nil` if the cinema hall is successfully added. Otherwise, it returns an error explaining the failure.
func (as *AdminService) AddNewCinemaHall(hallName, hallType string, capacity int) error {
	// Attempt to insert the new cinema hall into the database using the provided details.
	err := as.db.InsertNewCinemaHall(hallName, hallType, capacity)
	if err != nil {
		if errors.Is(err, models.ErrDuplicatedCinemaHall) {
			return ErrDuplicatedCinemaHall
		}
		// If an error occurs during insertion, return a wrapped error with context.
		return fmt.Errorf("error occurred while adding new hall: %w", err)
	}

	// Return nil if the cinema hall was successfully added.
	return nil
}

// FetchAllCinemaHalls retrieves all cinema halls from the database for the admin page.
//
// This function attempts to fetch a list of all cinema halls from the database for the admin interface.
// If an error occurs during the retrieval process, it returns a wrapped error. If successful, it returns the list of cinema halls.
//
// Parameters:
//   - None
//
// Returns:
//   - []models.CinemaHallForAdmin: A slice of cinema halls fetched from the database.
//   - error: Returns `nil` and the list of cinema halls if successful. Otherwise, it returns an error explaining the failure.
func (as *AdminService) FetchAllCinemaHalls() ([]models.CinemaHallForAdmin, error) {
	// Attempt to retrieve all cinema halls from the database.
	allCinemaHalls, err := as.db.RetrieveAllCinemaHallsForAdmin()
	if err != nil {
		// If the error is a "not found" error, return a custom error.
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return nil, ErrCinemaHallNotFound
		}
		// If another error occurs, return a wrapped error with context.
		return nil, fmt.Errorf("error occurred while fetching all cinema halls: %w", err)
	}

	// Return the list of cinema halls if successful.
	return allCinemaHalls, nil
}

// FetchCinemaHallInfo retrieves detailed information about a cinema hall from the database by its ID.
//
// This function attempts to fetch cinema hall information based on the provided `cinemaHallID` from the database.
// If an error occurs during the retrieval process, it returns an appropriate error. If successful, it returns the details of the cinema hall.
//
// Parameters:
//   - cinemaHallID (int): The unique identifier of the cinema hall to fetch information for.
//
// Returns:
//   - models.CinemaHallForAdmin: The details of the cinema hall fetched from the database.
//   - error: Returns `nil` and the cinema hall details if successful. Otherwise, it returns an error explaining the failure.
func (as *AdminService) FetchCinemaHallInfo(cinemaHallID int) (models.CinemaHallForAdmin, error) {
	// Attempt to retrieve the cinema hall information from the database.
	cinemaHall, err := as.db.RetrieveCinemaHallInfoByID(cinemaHallID)
	if err != nil {
		// If the error is a "not found" error, return a custom error for cinema hall not found.
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return models.CinemaHallForAdmin{}, ErrCinemaHallNotFound
		}
		// If another error occurs, return a wrapped error with additional context.
		return models.CinemaHallForAdmin{}, fmt.Errorf("error occurred while fetching cinema hall: %w", err)
	}

	// Return the retrieved cinema hall information if successful.
	return cinemaHall, nil
}

// DeleteCinemaHall removes a cinema hall from the database based on the provided cinema hall ID.
//
// This function attempts to delete the cinema hall from the database using the provided `cinemaHallID`.
// If the hall does not exist or an error occurs during the deletion, it returns an appropriate error message.
//
// Parameters:
//   - cinemaHallID (int): The unique identifier of the cinema hall to be deleted.
//
// Returns:
//   - error: Returns `nil` if the cinema hall was successfully deleted.
//     Otherwise, returns an error detailing the failure (e.g., if the hall was not found).
func (as *AdminService) DeleteCinemaHall(cinemaHallID int) error {
	// Attempt to delete the cinema hall by calling the database deletion function.
	err := as.db.DeleteCinemaHallByID(cinemaHallID)
	if err != nil {
		// If the error indicates the cinema hall was not found, return a custom error.
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return ErrCinemaHallNotFound
		}
		// If another error occurs, return a wrapped error with additional context.
		return fmt.Errorf("error occurred while deleting cinema hall: %w", err)
	}

	// Return `nil` if the deletion was successful, indicating no errors occurred.
	return nil
}

// AddCinemaSeats adds a new cinema seat to the database for a specific cinema hall.
//
// This function attempts to add a new seat with the provided details (seat row, seat number, seat type, hall ID) to the database.
// If the seat already exists or the hall is not found, it returns an appropriate error. If successful, it returns `nil`.
//
// Parameters:
//   - seatRow (string): The row identifier (e.g., "A", "B", "C") where the seat is located.
//   - seatNumber (int): The seat number within the given row (e.g., 1, 2, 3).
//   - seatType (string): The type of seat (e.g., "Standard", "VIP", etc.).
//   - hallID (int): The unique identifier of the cinema hall where the seat will be added.
//
// Returns:
//   - error: Returns `nil` if the seat was successfully added. Otherwise, it returns an error explaining the failure.
func (as *AdminService) AddCinemaSeats(seatRow string, seatNumber int, seatType string, hallID int) error {
	// Attempt to insert the new cinema seat into the database using the provided details.
	err := as.db.InsertCinemaSeats(seatRow, seatNumber, seatType, hallID)
	if err != nil {
		// If the error indicates that the seat already exists, return a specific error for that case.
		if errors.Is(err, models.ErrCinemaSeatAlreadyExists) {
			return ErrCinemaSeatAlreadyExists
		}

		// If the error indicates that the cinema hall is not found, return a custom error for hall not found.
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return ErrCinemaHallNotFound
		}

		// For other errors, wrap the error with context indicating where the failure occurred.
		return fmt.Errorf("error occurred while adding new cinema seat: %w", err)
	}

	// Return nil if the seat was successfully added.
	return nil
}

// FetchALLCinemaSeatsByHallID retrieves all cinema seats for a specific hall by its ID.
//
// This function fetches all the cinema seats associated with the given hall ID.
// If no seats are found or if there are any other errors during the process, it returns an appropriate error.
//
// Parameters:
//   - hallID (int): The unique identifier of the cinema hall for which the seats need to be fetched.
//
// Returns:
//   - []models.CinemaSeatForAdmin: A slice of CinemaSeatForAdmin models, representing all the seats in the given hall.
//   - error: Returns `nil` if successful, or an error explaining why the operation failed.
func (as *AdminService) FetchALLCinemaSeatsByHallID(hallID int) ([]models.CinemaSeatForAdmin, error) {

	// Attempt to retrieve all cinema seats for the given hall ID.
	allCinemaSeat, err := as.db.RetrieveALLCinemaSeatsByHallID(hallID)
	if err != nil {
		// If no seats are found (ErrCinemaSeatNotFound), return the specific error for that case.
		if errors.Is(err, models.ErrCinemaSeatNotFound) {
			return nil, ErrCinemaSeatNotFound
		}

		// For any other error, wrap the error with additional context to indicate where the error occurred.
		return nil, fmt.Errorf("error occurred while fetching all cinema seats by hall id: %w", err)
	}

	// Return the list of cinema seats if the operation is successful.
	return allCinemaSeat, nil
}

// DeleteCinemaSeat removes a cinema seat from the database by its unique ID.
//
// This function attempts to delete a cinema seat from the database using the given seat ID.
// If the seat is not found or there are other errors during the process, it returns an appropriate error.
//
// Parameters:
//   - cinemaSeatID (int): The unique identifier of the cinema seat to be deleted.
//
// Returns:
//   - error: Returns `nil` if the operation is successful, or an error explaining why the operation failed.
func (as *AdminService) DeleteCinemaSeat(cinemaSeatID int) error {

	// Attempt to delete the cinema seat by its unique ID.
	err := as.db.DeleteCinemaSeatByID(cinemaSeatID)
	if err != nil {
		// If the seat is not found (ErrCinemaSeatNotFound), return the specific error.
		if errors.Is(err, models.ErrCinemaSeatNotFound) {
			return ErrCinemaSeatNotFound
		}

		// For any other error, wrap the error with additional context indicating where the error occurred.
		return fmt.Errorf("error occurred while deleting cinema seat: %w", err)
	}

	// Return nil if the deletion is successful.
	return nil
}

// AddNewShow adds a new show to the database and assigns seats to it for the specified cinema hall.
//
// This function first formats the provided `showDate` and `startTime` as strings. Then it tries to insert a new show into
// the database and, if successful, proceeds to fetch all cinema seats for the specified hall. It inserts a new show seat
// for each available cinema seat in the hall, marking the status as "Available" and setting the price to 0.
//
// Parameters:
//   - showDate (time.Time): The date of the show in time format.
//   - startTime (time.Time): The start time of the show in time format.
//   - hallID (int): The ID of the cinema hall where the show will take place.
//   - movieID (int): The ID of the movie being shown.
//
// Returns:
//   - error: Returns `nil` if the operation is successful, or an error explaining why the operation failed.
func (as *AdminService) AddNewShow(showDate, startTime time.Time, hallID int, movieID int) error {

	// Format the provided date and time into strings for database insertion.
	formattedDate := showDate.Format("2006-01-02")
	formattedTime := startTime.Format("15:04:05")

	// Insert the new show into the database.
	showID, err := as.db.InsertNewShow(formattedDate, formattedTime, hallID, movieID)
	if err != nil {
		// : Handle errors related to missing cinema hall or movie.
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return ErrCinemaHallNotFound
		}
		if errors.Is(err, models.ErrMovieNotFoundByID) {
			return ErrMovieNotFoundByID
		}

		if errors.Is(err, models.ErrShowAlreadyExists) {
			return ErrShowAlreadyExists
		}

		// : Return any other errors that occurred during the insertion of the new show.
		return fmt.Errorf("error occurred while adding new show: %w", err)
	}

	// Retrieve all cinema seats for the specified cinema hall.
	allCinemaSeats, err := as.db.RetrieveALLCinemaSeatsByHallID(hallID)
	if err != nil {
		// : Handle the error if no cinema seats are found for the hall.
		if errors.Is(err, models.ErrCinemaSeatNotFound) {
			return ErrCinemaSeatNotFound
		}

		// : Return any other errors that occurred during the retrieval of cinema seats.
		return fmt.Errorf("error occurred while retrieving cinema hall seats: %w", err)
	}

	// Insert a new show seat for each cinema seat in the hall, with initial status "Available" and price 0.
	for _, cinemaSeat := range allCinemaSeats {
		err = as.db.InsertNewShowSeat("Available", 0, cinemaSeat.CinemaSeatID, showID)
		if err != nil {
			// : Handle the error if a show seat cannot be found or inserted.
			if errors.Is(err, models.ErrShowSeatNotFound) {
				return ErrShowSeatNotFound
			}

			// : Return any other errors that occurred during the insertion of show seats.
			return fmt.Errorf("error occurred while inserting new show seats: %w", err)
		}
	}

	// Return nil indicating that the new show and its seats were successfully added.
	return nil
}

// FetchAllShowsForAdmin retrieves all shows from the database for the admin page.
//
// This function retrieves a list of all shows in the system from the database.
// If no shows are found or any other error occurs, it returns an appropriate error.
//
// Parameters:
//   - None
//
// Returns:
//   - ([]models.ShowForAdmin): A list of shows for the admin to manage.
//   - error: Returns `nil` if the operation is successful, or an error explaining why the operation failed.
func (as *AdminService) FetchAllShowsForAdmin() ([]models.ShowForAdmin, error) {

	// Attempt to retrieve all shows for the admin from the database.
	allShows, err := as.db.RetrieveAllShowsForAdmin()
	if err != nil {
		// If no shows are found (ErrShowNotFound), return the specific error.
		if errors.Is(err, models.ErrShowNotFound) {
			return nil, ErrShowNotFound
		}
		// If any other error occurs, wrap the error with additional context.
		return nil, fmt.Errorf("error occurred while fetching all shows: %w", err)
	}

	// Return the list of shows if retrieval is successful.
	return allShows, nil
}

// UpdateShow updates the details of an existing show in the database based on the provided showID.
//
// This function formats the provided `showDate` and `startTime` as strings, then attempts to update the
// show details (such as date, time, hall, and movie) in the database. If any errors occur during the process,
// appropriate error messages are returned.
//
// Parameters:
//   - showID (int): The unique identifier for the show to be updated.
//   - showDate (time.Time): The new date of the show.
//   - startTime (time.Time): The new start time of the show.
//   - hallID (int): The ID of the cinema hall where the show will take place.
//   - movieID (int): The ID of the movie being shown.
//
// Returns:
//   - error: Returns `nil` if the update operation is successful, or an error explaining why the operation failed.
func (as *AdminService) UpdateShow(showID int, showDate, startTime time.Time, hallID int, movieID int) error {

	// Format the show date and start time to match the database format.
	formattedDate := showDate.Format("2006-01-02")
	formattedTime := startTime.Format("15:04:05")

	// Attempt to update the show details in the database.
	err := as.db.UpdateShowByID(showID, formattedDate, formattedTime, hallID, movieID)
	if err != nil {
		// : Check if the error is due to a non-existing cinema hall.
		if errors.Is(err, models.ErrCinemaHallNotFound) {
			return ErrCinemaHallNotFound
		}
		// : Check if the error is due to a non-existing movie.
		if errors.Is(err, models.ErrMovieNotFoundByID) {
			return ErrMovieNotFoundByID
		}
		// : Check if the error is due to a non-existing show.
		if errors.Is(err, models.ErrShowNotFound) {
			return ErrShowNotFound
		}
		// : Return any other error that might have occurred.
		return fmt.Errorf("error occurred while updating show data: %w", err)
	}

	// Return nil if the show update was successful.
	return nil
}

// DeleteShow deletes a show from the database based on the provided showID.
//
// This function attempts to delete the show from the database. If the show does not exist,
// it returns an appropriate error message. If any other errors occur during the process,
// they are returned as wrapped errors.
//
// Parameters:
//   - showID (int): The unique identifier for the show to be deleted.
//
// Returns:
//   - error: Returns `nil` if the delete operation is successful, or an error explaining why the operation failed.
func (as *AdminService) DeleteShow(showID int) error {
	// Attempt to delete the show by its ID.
	err := as.db.DeleteShowByID(showID)
	if err != nil {
		// : Check if the error is due to a non-existing show.
		if errors.Is(err, models.ErrShowNotFound) {
			return ErrShowNotFound
		}
		// : Return any other error that might have occurred during the delete operation.
		return fmt.Errorf("error occurred while deleting show: %w", err)
	}

	// Return nil if the show was successfully deleted.
	return nil
}

// FetchAllShowSeats retrieves all show seats for a specific show from the database based on the provided showID.
//
// This function fetches all the seats associated with a given show. If no seats are found or if there
// is an error during the retrieval, an appropriate error is returned. If the retrieval is successful,
// a list of show seats is returned.
//
// Parameters:
//   - showID (int): The unique identifier for the show whose seats need to be fetched.
//
// Returns:
//   - ([]models.ShowSeatForAdmin, error): A slice of ShowSeatForAdmin models containing all the show seats,
//     or an error if the retrieval fails.
func (as *AdminService) FetchAllShowSeats(showID int) ([]models.ShowSeatForAdmin, error) {

	// Attempt to retrieve all seats for the given show.
	allShowSeats, err := as.db.RetrieveAllShowSeats(showID)
	if err != nil {
		// : If no seats are found for the show, return an appropriate error.
		if errors.Is(err, models.ErrShowSeatNotFound) {
			return nil, ErrShowSeatNotFound
		}
		// : Return any other error encountered during the seat retrieval process.
		return nil, fmt.Errorf("error occurred while fetching all show seats: %w", err)
	}

	for i := range allShowSeats {
		// Convert integer seat price (in cents) to float32 and divide by 100 to display as a price
		allShowSeats[i].SeatPrice = float32(allShowSeats[i].SeatPrice) / 100.0
	}

	// Return the list of show seats if the retrieval was successful.
	return allShowSeats, nil
}

// UpdateShowSeat updates the price of a specific show seat identified by showSeatID.
//
// This function updates the price of a seat associated with a specific show. It takes the new price as
// a float (multiplied by 10) and the unique identifier for the show seat. If the seat cannot be found or
// there is an error during the update, an appropriate error is returned. If successful, the function returns nil.
//
// Parameters:
//   - seatPrice (float32): The new price of the seat to be updated (in the form of float32, multiplied by 10).
//   - showSeatID (int): The unique identifier for the show seat whose price needs to be updated.
//
// Returns:
//   - error: Returns nil if the update was successful, or an error if the update fails.
func (as *AdminService) UpdateShowSeat(seatPrice float32, showSeatID int) error {

	// Adjust the seat price (multiply by 10 to match the expected data format).
	seatPriceInCents := int(seatPrice * 100)

	// Attempt to update the show seat's price in the database.
	err := as.db.UpdateShowSeatByID(seatPriceInCents, showSeatID)
	if err != nil {
		// : If no show seat is found for the given showSeatID, return a specific error.
		if errors.Is(err, models.ErrShowSeatNotFound) {
			return ErrShowSeatNotFound
		}
		// : Return any other error that occurs during the update process.
		return fmt.Errorf("error occurred while updating seat price: %w", err)
	}

	// Return nil if the seat price update was successful.
	return nil
}
