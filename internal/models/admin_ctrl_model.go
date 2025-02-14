package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type DBContractAdminCtrl interface {
	InsertNewCarouselImages(imageURL, title, description string, orderPriority int) error
	RetrieveCarouselImagesForAdmin() ([]CarouselImageForAdmin, error)
	UpdateCarouselImagesByID(imageURL, title, description string, orderPriority int, carouselImageID int) error
	DeleteCarouselImagesByID(carouselImageID int) error

	InsertNewMovie(title string, description, genre, language, trailerURL, posterURL string, rating int, ratingProvider string, duration int, releaseDate, ageLimit string) error
	RetrieveAllMoviesForAdmin() ([]AllMoviesForAdmin, error)
	RetrieveAMovieForAdmin(movieID int) (MovieForAdmin, error)
	UpdateMovieInfoForAdminByMovieID(movieID int, title, description, genre, language, trailerURL, posterURL string, rating float32, ratingProvider string, duration int, relaseDate string, ageLimit string) error
	DeleteMovieByMovieID(movieID int) error

	InsertActorsCrew(fullName, imageURL, occupation, roleDescription, bornDate, birthplace, about string, isActor bool) (int, error)
	InsertMovieActorCrew(movieID, actorCrewID int) error
	RetrieveAllActorsCrewByMovieID(movieID int) ([]AllActorCrewForAdmin, error)
	RetrieveActorCrew(actorCrewID int) (ActorCrewForAdmin, error)
	UpdateActorCrewInformation(fullName, imageURL, occupation, roleDescription, bornDate, birthplace, about string, isActor bool, actorCrewID int) error
	DeleteActorCrewByID(actorCrewID int) error

	InsertNewCinemaHall(hallName, hallType string, capacity int) error
	RetrieveAllCinemaHallsForAdmin() ([]CinemaHallForAdmin, error)
	RetrieveCinemaHallInfoByID(cinemaHallID int) (CinemaHallForAdmin, error)

	DeleteCinemaHallByID(cinemaHallID int) error

	CountCinemaSeatsByHallID(hallID int) (int, error)
	InsertCinemaSeats(seatRow string, seatNumber int, seatType string, hallID int) error
	RetrieveALLCinemaSeatsByHallID(hallID int) ([]CinemaSeatForAdmin, error)
	DeleteCinemaSeatByID(cinemaSeatID int) error

	InsertNewShow(showDate string, startTime string, hallID int, movieID int) (int, error)
	RetrieveAllShowsForAdmin() ([]ShowForAdmin, error)
	UpdateShowByID(showID int, showDate string, startTime string, hallID int, movieID int) error
	DeleteShowByID(showID int) error

	InsertNewShowSeat(seatStatus string, seatPrice int, cinemSeatID int, showID int) error
	RetrieveAllShowSeats(showID int) ([]ShowSeatForAdmin, error)
	UpdateShowSeatByID(seatPrice int, showSeatID int) error
}

type AdminOperations struct {
	db DBContractAdminCtrl
}

func NewAdminOperations(db DBContractAdminCtrl) *AdminOperations {
	return &AdminOperations{db: db}
}

// InsertNewCarouselImages inserts a new carousel image record into the database.
//
// Parameters:
//   - imageURL (string): The URL of the image to be inserted.
//   - title (string): The title associated with the image.
//   - description (string): The description of the image.
//   - orderPriority (int): The order priority for the carousel image (used to determine the order of display).
//
// Returns:
//   - error: Returns nil if the insertion is successful. If there's an error during execution, it returns a wrapped error providing context for the failure.
func (psql *Postgres) InsertNewCarouselImages(imageURL, title, description string, orderPriority int) error {
	// SQL statement for inserting a new carousel image record
	stmt := `INSERT INTO carousel_images (image_url, title, description, order_priority) VALUES ($1, $2, $3, $4)`

	// Execute the SQL statement with the provided parameters
	_, err := psql.DB.Exec(stmt, imageURL, title, description, orderPriority)
	if err != nil {
		// Return a wrapped error with context if insertion fails
		return fmt.Errorf("failed to insert data into carousel images: %w", err)
	}

	// Return nil if the insertion was successful
	return nil
}

// RetrieveCarouselImagesForAdmin retrieves all carousel images from the database for display on the admin page.
//
// Parameters:
//   - None: This function doesn't take any parameters.
//
// Returns:
//   - []CarouselImageForAdmin: A slice of CarouselImageForAdmin structures containing the carousel images fetched from the database.
//   - error: Returns nil if the retrieval is successful. If an error occurs during the query execution or scanning rows, an error is returned providing context.
func (psql *Postgres) RetrieveCarouselImagesForAdmin() ([]CarouselImageForAdmin, error) {
	// SQL query to select all carousel images from the database
	stmt := `SELECT * FROM carousel_images`

	// Execute the query
	rows, err := psql.DB.Query(stmt)
	if err != nil {
		// Return a wrapped error if query execution fails
		return nil, fmt.Errorf("failed to retrieve carousel images for admin page: %w", err)
	}

	// Ensure rows are closed after the function completes
	defer rows.Close()

	// Initialize a slice to hold the carousel images
	var carouselImages []CarouselImageForAdmin

	// Iterate through the rows returned by the query
	for rows.Next() {
		var carouselImage CarouselImageForAdmin

		// Scan the row data into the carouselImage struct
		err := rows.Scan(&carouselImage.CarouselImageID, &carouselImage.ImageURL, &carouselImage.Title, &carouselImage.Description, &carouselImage.OrderPriority, &carouselImage.CreatedAt, &carouselImage.UpdatedAt)
		if err != nil {
			// Check for the "no rows" error and return a specific error if no images are found
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrAdminPageCarouselImagesNotFound
			}
			// Return a wrapped error if row scanning fails
			return nil, fmt.Errorf("failed to retrieve admin page carousel images: %w", err)
		}

		// Append the successfully scanned image to the slice
		carouselImages = append(carouselImages, carouselImage)
	}

	// Return the slice of carousel images
	return carouselImages, nil
}

// UpdateCarouselImagesByID updates a carousel image's details (URL, title, description, and order priority)
// in the database based on the provided image ID.
//
// Parameters:
//   - imageURL (string): The new URL of the image to update.
//   - title (string): The new title of the image.
//   - description (string): The new description of the image.
//   - orderPriority (int): The new order priority for the carousel image.
//   - carouselImageID (int): The unique ID of the carousel image to update.
//
// Returns:
//   - error: Returns nil if the update is successful. If an error occurs during execution (such as failure to execute the query
//     or check rows affected), an appropriate error is returned. If no rows were affected (i.e., the image ID wasn't found),
//     the function returns `ErrAdminPageCarouselImagesNotFound`.
func (psql *Postgres) UpdateCarouselImagesByID(imageURL, title, description string, orderPriority int, carouselImageID int) error {
	// SQL query to update the carousel image details by ID
	stmt := `UPDATE carousel_images SET image_url = $1, title = $2, description = $3, order_priority = $4 WHERE id = $5`

	// Execute the update statement with the provided parameters
	result, err := psql.DB.Exec(stmt, imageURL, title, description, orderPriority, carouselImageID)
	if err != nil {
		// Return a wrapped error if query execution fails
		return fmt.Errorf("failed to update carousel image: %w", err)
	}

	// Check how many rows were affected by the update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Return an error if there's an issue checking rows affected
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// If no rows were affected, the image ID was not found
	if rowsAffected == 0 {
		return ErrAdminPageCarouselImagesNotFound
	}

	// Return nil if the update was successful
	return nil
}

// DeleteCarouselImagesByID deletes a carousel image from the database based on the provided image ID.
//
// Parameters:
//   - carouselImageID (int): The unique ID of the carousel image to delete.
//
// Returns:
//   - error: Returns nil if the deletion is successful. If an error occurs during the query execution,
//     or when checking the rows affected, an error is returned with context. If no rows were affected (i.e.,
//     no image with the provided ID exists), the function returns `ErrAdminPageCarouselImagesNotFound`.
func (psql *Postgres) DeleteCarouselImagesByID(carouselImageID int) error {
	// SQL query to delete the carousel image by ID
	stmt := `DELETE FROM carousel_images WHERE id = $1`

	// Execute the delete statement with the provided image ID
	result, err := psql.DB.Exec(stmt, carouselImageID)
	if err != nil {
		// Return a wrapped error if query execution fails
		return fmt.Errorf("failed to delete carousel image by id: %w", err)
	}

	// Check how many rows were affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Return an error if there's an issue checking rows affected
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// If no rows were affected, the image ID was not found
	if rowsAffected == 0 {
		return ErrAdminPageCarouselImagesNotFound
	}

	// Return nil if the deletion was successful
	return nil
}

// InsertNewMovie inserts a new movie record into the database.
//
// Parameters:
//   - title (string): The title of the movie.
//   - description (string): A brief description of the movie.
//   - genre (string): The genre of the movie (e.g., Action, Comedy, Drama).
//   - language (string): The language in which the movie is available.
//   - trailerURL (string): The URL of the movie's trailer.
//   - posterURL (string): The URL of the movie's poster image.
//   - rating (int): The rating of the movie (e.g., IMDB rating, out of 10).
//   - ratingProvider (string): The source or provider of the movie rating (e.g., IMDB, Rotten Tomatoes).
//   - duration (int): The duration of the movie in minutes.
//   - releaseDate (string): The release date of the movie (in string format, e.g., "YYYY-MM-DD").
//   - ageLimit (string): The age limit or rating (e.g., PG-13, R).
//
// Returns:
//   - error: Returns nil if the insertion is successful. If an error occurs during execution (e.g., query execution fails),
//     it returns a wrapped error with context.
func (psql *Postgres) InsertNewMovie(title string, description, genre, language, trailerURL, posterURL string, rating int, ratingProvider string, duration int, releaseDate, ageLimit string) error {
	// SQL statement to insert a new movie record into the database
	stmt := `INSERT INTO movies (title, description, genre, language, trailer_url, poster_url, rating, rating_provider, duration, release_date, age_limit) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	// Execute the SQL statement with the provided movie details
	_, err := psql.DB.Exec(stmt, title, description, genre, language, trailerURL, posterURL, rating, ratingProvider, duration, releaseDate, ageLimit)
	if err != nil {
		// Return a wrapped error if query execution fails
		return fmt.Errorf("failed to insert data into movies: %w", err)
	}

	// Return nil if the insertion was successful
	return nil
}

// RetrieveAllMoviesForAdmin retrieves all movie records from the database, specifically for the admin page.
//
// Parameters:
//   - None: This function doesn't take any parameters.
//
// Returns:
//   - []AllMoviesForAdmin: A slice of `AllMoviesForAdmin` structs containing the movie details (ID and title) for the admin page.
//   - error: Returns nil if the retrieval is successful. If an error occurs during the query execution, or while scanning rows,
//     an error is returned with context. If no rows are found, it returns `ErrAdminPageMovieNotFound`.
func (psql *Postgres) RetrieveAllMoviesForAdmin() ([]AllMoviesForAdmin, error) {
	// SQL query to select the ID and title of all movies from the database
	stmt := `SELECT id, title FROM movies`

	// Execute the query to retrieve movie data
	rows, err := psql.DB.Query(stmt)
	if err != nil {
		// Return a wrapped error if the query execution fails
		return nil, fmt.Errorf("failed to retrieve all movies for admin page: %w", err)
	}

	// Ensure rows are closed after the function completes
	defer rows.Close()

	// Initialize a slice to hold the movie data for the admin page
	var allMovies []AllMoviesForAdmin

	// Iterate through the rows returned by the query
	for rows.Next() {
		var movie AllMoviesForAdmin

		// Scan the row data into the movie struct
		err := rows.Scan(&movie.ID, &movie.Title)
		if err != nil {
			// Check for the "no rows" error and return a specific error if no movies are found
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrAdminPageMovieNotFound
			}
			// Return a wrapped error if scanning the rows fails
			return nil, fmt.Errorf("failed to retrieve all movies for admin page: %w", err)
		}

		// Append the successfully scanned movie to the slice
		allMovies = append(allMovies, movie)
	}

	// Return the slice of movie data
	return allMovies, nil
}

// RetrieveAMovieForAdmin retrieves a specific movie record from the database based on the provided movie ID,
// specifically for the admin page.
//
// Parameters:
//   - movieID (int): The unique ID of the movie to retrieve.
//
// Returns:
//   - MovieForAdmin: A struct containing the details of the requested movie (MovieForAdmin).
//   - error: Returns nil if the retrieval is successful. If an error occurs during query execution or scanning,
//     an error is returned with context. If no movie with the specified ID is found, it returns `ErrAdminPageMovieNotFound`.
func (psql *Postgres) RetrieveAMovieForAdmin(movieID int) (MovieForAdmin, error) {
	// SQL query to select all columns for a movie by its ID
	stmt := `SELECT * FROM movies WHERE id = $1`

	// Declare a variable to hold the movie data
	var movie MovieForAdmin

	// Execute the query and scan the result into the movie variable
	err := psql.DB.QueryRow(stmt, movieID).Scan(&movie.MovieID, &movie.Title, &movie.Description, &movie.Genere, &movie.Language, &movie.TrailerURL, &movie.PosterURL, &movie.Rating, &movie.RatingProvider, &movie.Duration, &movie.RelaseDate, &movie.AgeLimit, &movie.CreatedAt, &movie.UpdatesAt)
	if err != nil {
		// If no rows are returned (i.e., the movie ID doesn't exist), return a specific error
		if errors.Is(err, sql.ErrNoRows) {
			return MovieForAdmin{}, ErrAdminPageMovieNotFound
		}
		// Return a wrapped error if scanning the result fails
		return MovieForAdmin{}, fmt.Errorf("failed to retrieve a movie for admin page: %w", err)
	}

	// Return the movie struct if retrieval was successful
	return movie, nil
}

// UpdateMovieInfoForAdminByMovieID updates the information of a specific movie in the database based on the provided movie ID.
//
// Parameters:
//   - movieID (int): The unique ID of the movie to update.
//   - title (string): The new title of the movie.
//   - description (string): The new description of the movie.
//   - genre (string): The new genre of the movie (e.g., Action, Comedy, Drama).
//   - language (string): The new language in which the movie is available.
//   - trailerURL (string): The new URL of the movie's trailer.
//   - posterURL (string): The new URL of the movie's poster image.
//   - rating (float32): The new rating of the movie (e.g., IMDB rating, out of 10).
//   - ratingProvider (string): The new source or provider of the movie rating (e.g., IMDB, Rotten Tomatoes).
//   - duration (int): The new duration of the movie in minutes.
//   - releaseDate (string): The new release date of the movie (in string format, e.g., "YYYY-MM-DD").
//   - ageLimit (string): The new age limit or rating (e.g., PG-13, R).
//
// Returns:
//   - error: Returns nil if the update is successful. If an error occurs during the query execution, or while checking rows affected,
//     an error is returned with context. If no rows are affected (i.e., the movie ID is not found), it returns `ErrAdminPageMovieNotFound`.
func (psql *Postgres) UpdateMovieInfoForAdminByMovieID(movieID int, title, description, genre, language, trailerURL, posterURL string, rating float32, ratingProvider string, duration int, relaseDate string, ageLimit string) error {
	// SQL query to update movie information based on movie ID
	stmt := `UPDATE movies SET title = $1, description = $2, genre = $3, language = $4, trailer_url = $5, poster_url = $6, rating = $7, rating_provider = $8, duration = $9, release_date = $10, age_limit = $11, updated_at = CURRENT_TIMESTAMP WHERE id = $12;`

	// Execute the SQL query with the provided parameters
	result, err := psql.DB.Exec(stmt, title, description, genre, language, trailerURL, posterURL, rating, ratingProvider, duration, relaseDate, ageLimit, movieID)
	if err != nil {
		// Return a wrapped error if query execution fails
		return fmt.Errorf("failed to update movie information: %w", err)
	}

	// Check how many rows were affected by the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Return an error if there's an issue checking rows affected
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// If no rows were affected, the movie ID was not found
	if rowsAffected == 0 {
		return ErrAdminPageMovieNotFound
	}

	// Return nil if the update was successful
	return nil
}

// DeleteMovieByMovieID deletes a movie record from the database based on the provided movie ID.
//
// Parameters:
//   - movieID (int): The unique ID of the movie to delete.
//
// Returns:
//   - error: Returns nil if the deletion is successful. If an error occurs during the query execution, or while checking rows affected,
//     an error is returned with context. If no rows are affected (i.e., no movie with the specified ID exists),
//     it returns `ErrAdminPageMovieNotFound`.
func (psql *Postgres) DeleteMovieByMovieID(movieID int) error {
	// SQL query to delete the movie by its ID
	stmt := `DELETE FROM movies WHERE id = $1;`

	// Execute the query to delete the movie with the provided ID
	result, err := psql.DB.Exec(stmt, movieID)
	if err != nil {
		// Return a wrapped error if query execution fails
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	// Check how many rows were affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Return an error if there's an issue checking rows affected
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// If no rows were affected, the movie ID was not found
	if rowsAffected == 0 {
		return ErrAdminPageMovieNotFound
	}

	// Return nil if the deletion was successful
	return nil
}

// InsertActorsCrew inserts a new actor/crew member record into the database and returns the generated ID.
//
// Parameters:
//   - fullName (string): The full name of the actor or crew member.
//   - imageURL (string): The URL to the image of the actor or crew member.
//   - occupation (string): The occupation of the individual (e.g., "Actor", "Director").
//   - roleDescription (string): A description of the individual's role in the project.
//   - bornDate (string): The birth date of the individual (in string format, e.g., "YYYY-MM-DD").
//   - birthplace (string): The birthplace of the individual.
//   - about (string): A brief biography or description of the individual.
//   - isActor (bool): A boolean indicating whether the person is an actor (true) or part of the crew (false).
//
// Returns:
//   - int: The ID of the newly inserted actor/crew member.
//   - error: Returns nil if the insertion is successful. If an error occurs during the query execution or scanning the returned ID,
//     it returns a wrapped error with context.
func (psql *Postgres) InsertActorsCrew(fullName, imageURL, occupation, roleDescription, bornDate, birthplace, about string, isActor bool) (int, error) {
	// SQL query to insert actor/crew information and return the generated ID
	stmt := `INSERT INTO actors_crew (full_name, image_url, occupation, role_description, born_date, birthplace, about, is_actor) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	// Variable to hold the generated ID of the newly inserted actor/crew member
	var actorCrewID int

	// Execute the query and scan the generated ID into actorCrewID
	err := psql.DB.QueryRow(stmt, fullName, imageURL, occupation, roleDescription, bornDate, birthplace, about, isActor).Scan(&actorCrewID)
	if err != nil {
		// Return a wrapped error if query execution or scanning the ID fails
		return 0, fmt.Errorf("failed to insert actors_crew info into database: %w", err)
	}

	// Return the generated ID and nil if the insertion was successful
	return actorCrewID, nil
}

// InsertMovieActorCrew inserts a new record into the movie_actors_crew table,
// associating a movie with an actor or crew member.
//
// Parameters:
//   - movieID (int): The unique ID of the movie to associate with the actor/crew member.
//   - actorCrewID (int): The unique ID of the actor or crew member to associate with the movie.
//
// Returns:
//   - error: Returns nil if the insertion is successful. If an error occurs during the query execution,
//     it returns a wrapped error with context.
func (psql *Postgres) InsertMovieActorCrew(movieID, actorCrewID int) error {
	// SQL query to insert the association between the movie and actor/crew member
	stmt := `INSERT INTO movie_actors_crew (movie_id, actor_crew_id) VALUES ($1, $2)`

	// Execute the query to insert the movie and actor/crew member association
	_, err := psql.DB.Exec(stmt, movieID, actorCrewID)
	if err != nil {
		// Return a wrapped error if the query execution fails
		return fmt.Errorf("failed to insert into movie_actors_crew: %w", err)
	}

	// Return nil if the insertion was successful
	return nil
}

// RetrieveAllActorsCrewByMovieID retrieves all actors and crew members associated with a given movie for the admin page.
//
// Parameters:
//   - movieID (int): The unique ID of the movie for which the actors and crew are being retrieved.
//
// Returns:
//   - []AllActorCrewForAdmin: A slice containing all the actor/crew members associated with the movie.
//   - error: Returns nil and the list of actor/crew members if successful. If an error occurs during the query execution,
//     or while scanning the results, it returns a wrapped error. If no actors/crew members are found for the given movie,
//     it returns `ErrActorCrewNotFound`.
func (psql *Postgres) RetrieveAllActorsCrewByMovieID(movieID int) ([]AllActorCrewForAdmin, error) {
	// SQL query to retrieve all actors and crew members associated with the specified movie ID
	stmt := `SELECT ac.id, ac.full_name, ac.image_url, ac.role_description, ac.is_actor FROM actors_crew ac 
             JOIN movie_actors_crew mac ON ac.id = mac.actor_crew_id 
             JOIN movies m ON mac.movie_id = m.id 
             WHERE m.id = $1`

	// Execute the query to retrieve the rows associated with the given movie ID
	rows, err := psql.DB.Query(stmt, movieID)
	if err != nil {
		// Return a wrapped error if the query fails
		return nil, fmt.Errorf("failed to retrieve all actors/crew by movieID for admin page: %w", err)
	}

	// Ensure rows are properly closed after processing
	defer rows.Close()

	// Slice to hold the resulting actor/crew members
	var allActorCrew []AllActorCrewForAdmin

	// Loop through each row and scan the result into a struct
	for rows.Next() {
		var actorCrew AllActorCrewForAdmin

		// Scan the row into the actorCrew struct
		err := rows.Scan(&actorCrew.ActorCrewID, &actorCrew.FullName, &actorCrew.ImageURL, &actorCrew.RoleDescription, &actorCrew.IsActor)
		if err != nil {
			// If no rows are found, return a specific error
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrActorCrewNotFound
			}
			// Return a wrapped error if scanning the row fails
			return nil, fmt.Errorf("failed to scan actors/crew by movieID for admin page: %w", err)
		}

		// Append the actor/crew member to the result slice
		allActorCrew = append(allActorCrew, actorCrew)
	}

	// Return the slice of actor/crew members
	return allActorCrew, nil
}

// RetrieveActorCrew retrieves the details of an actor or crew member based on their unique ID.
//
// Parameters:
//   - actorCrewID (int): The unique ID of the actor or crew member to retrieve.
//
// Returns:
//   - ActorCrewForAdmin: A struct containing the actor or crew member's details.
//   - error: Returns nil if the retrieval is successful. If an error occurs during the query execution or scanning the result,
//     it returns a wrapped error with context. If no record is found with the given ID, it returns `ErrActorCrewNotFound`.
func (psql *Postgres) RetrieveActorCrew(actorCrewID int) (ActorCrewForAdmin, error) {
	// SQL query to retrieve actor/crew member details by their ID
	stmt := `SELECT * FROM actors_crew WHERE id = $1`

	// Variable to store the result of the query
	var actorCrew ActorCrewForAdmin

	// Execute the query and scan the result into actorCrew
	err := psql.DB.QueryRow(stmt, actorCrewID).Scan(&actorCrew.ActorCrewID, &actorCrew.FullName, &actorCrew.ImageURL,
		&actorCrew.Occupation, &actorCrew.RoleDescription, &actorCrew.BornDate, &actorCrew.Birthplace,
		&actorCrew.About, &actorCrew.IsActor)

	// Check for errors during the scan process
	if err != nil {
		// If no rows are found for the given ID, return a specific error
		if errors.Is(err, sql.ErrNoRows) {
			return ActorCrewForAdmin{}, ErrActorCrewNotFound
		}
		// Return a wrapped error if scanning the row fails
		return ActorCrewForAdmin{}, fmt.Errorf("failed to scan actorCrew for admin page: %w", err)
	}

	// Return the populated actorCrew struct if successful
	return actorCrew, nil
}

// UpdateActorCrewInformation updates the details of an existing actor or crew member based on their ID.
//
// Parameters:
//   - fullName (string): The full name of the actor or crew member.
//   - imageURL (string): The URL of the actor/crew member's image.
//   - occupation (string): The occupation of the actor or crew member (e.g., actor, director).
//   - roleDescription (string): A description of the role played by the actor or the responsibilities of the crew member.
//   - bornDate (string): The birth date of the actor/crew member (in string format).
//   - birthplace (string): The birthplace of the actor/crew member.
//   - about (string): Additional information about the actor/crew member.
//   - isActor (bool): A flag indicating whether the person is an actor or a crew member.
//   - actorCrewID (int): The unique ID of the actor or crew member to update.
//
// Returns:
//   - error: Returns nil if the update is successful. If an error occurs during the query execution or checking rows affected,
//     it returns a wrapped error with context. If no rows are affected (i.e., the record is not found), it returns `ErrActorCrewNotFound`.
func (psql *Postgres) UpdateActorCrewInformation(fullName, imageURL, occupation, roleDescription, bornDate, birthplace, about string, isActor bool, actorCrewID int) error {

	// SQL query to update the actor/crew member's information based on the given ID
	stmt := `UPDATE actors_crew SET full_name = $1, image_url = $2, occupation = $3, role_description = $4, born_date = $5, birthplace = $6, about = $7, is_actor = $8 WHERE id = $9`

	// Execute the query to update the actor/crew member's details in the database
	result, err := psql.DB.Exec(stmt, fullName, imageURL, occupation, roleDescription, bornDate, birthplace, about, isActor, actorCrewID)
	if err != nil {
		// Return a wrapped error if the query execution fails
		return fmt.Errorf("failed to update actorCrew information: %w", err)
	}

	// Check the number of rows affected by the update operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Return a wrapped error if there's an issue checking rows affected
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// If no rows are affected, it means no matching record was found, so return a specific error
	if rowsAffected == 0 {
		return ErrActorCrewNotFound
	}

	// Return nil if the update was successful
	return nil
}

// DeleteActorCrewByID deletes an actor or crew member from the database based on their unique ID.
//
// Parameters:
//   - actorCrewID (int): The unique ID of the actor or crew member to delete.
//
// Returns:
//   - error: Returns nil if the deletion is successful. If an error occurs during the query execution or checking rows affected,
//     it returns a wrapped error with context. If no rows are affected (i.e., the record is not found), it returns `ErrActorCrewNotFound`.
func (psql *Postgres) DeleteActorCrewByID(actorCrewID int) error {
	// SQL query to delete an actor/crew member by their ID
	stmt := `DELETE FROM actors_crew WHERE id = $1;`

	// Execute the delete query and capture the result
	result, err := psql.DB.Exec(stmt, actorCrewID)
	if err != nil {
		// Return a wrapped error if the query execution fails
		return fmt.Errorf("failed to delete actorCrew: %w", err)
	}

	// Check how many rows were affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Return a wrapped error if there's an issue checking rows affected
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// If no rows were affected, it means no matching record was found, return a specific error
	if rowsAffected == 0 {
		return ErrActorCrewNotFound
	}

	// Return nil if the deletion was successful
	return nil
}

// InsertNewCinemaHall inserts a new cinema hall into the database.
//
// Parameters:
//   - hallName (string): The name of the cinema hall.
//   - hallType (string): The type of the cinema hall (e.g., IMAX, regular, etc.).
//   - capacity (int): The seating capacity of the cinema hall.
//
// Returns:
//   - error: Returns nil if the insertion is successful. If an error occurs during the query execution,
//     it returns a wrapped error with context.
func (psql *Postgres) InsertNewCinemaHall(hallName, hallType string, capacity int) error {
	// SQL query to insert a new cinema hall into the database
	stmt := `INSERT INTO cinema_hall (hall_name, hall_type, capacity) VALUES ($1, $2, $3)`

	// Execute the query to insert the new cinema hall
	_, err := psql.DB.Exec(stmt, hallName, hallType, capacity)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			// Handle duplicate entry error gracefully
			return ErrDuplicatedCinemaHall
		}

		// Return a wrapped error if the query execution fails
		return fmt.Errorf("failed to insert new hall information: %w", err)
	}

	// Return nil if the insertion is successful
	return nil
}

// RetrieveAllCinemaHallsForAdmin retrieves all cinema hall information for the admin page.
//
// Returns:
//   - []CinemaHallForAdmin: A slice of CinemaHallForAdmin containing information of all cinema halls.
//   - error: Returns an error if there's a problem retrieving the cinema hall data, scanning the rows, or no halls are found.
func (psql *Postgres) RetrieveAllCinemaHallsForAdmin() ([]CinemaHallForAdmin, error) {
	// SQL query to retrieve all cinema hall records from the database
	stmt := `SELECT * FROM cinema_hall`

	// Execute the query to fetch rows
	rows, err := psql.DB.Query(stmt)
	if err != nil {
		// Return a wrapped error if the query execution fails
		return nil, fmt.Errorf("failed to retrieve all hall information for admin page: %w", err)
	}

	// Ensure rows are closed after the function finishes
	defer rows.Close()

	// Slice to hold all the retrieved cinema hall records
	var allCinemaHalls []CinemaHallForAdmin

	// Iterate over the rows and scan each row into the cinemaHall struct
	for rows.Next() {
		var cinemaHall CinemaHallForAdmin

		// Scan the columns of the current row into the cinemaHall struct
		err := rows.Scan(&cinemaHall.CinemaHallID, &cinemaHall.HallName, &cinemaHall.HallType, &cinemaHall.Capacity)
		if err != nil {
			// Check if no rows were found and return a custom error if so
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrCinemaHallNotFound
			}
			// Return a wrapped error if there's an issue scanning the row
			return nil, fmt.Errorf("failed to scan cinema hall information: %w", err)
		}

		// Append the scanned cinema hall to the slice
		allCinemaHalls = append(allCinemaHalls, cinemaHall)
	}

	// Return the slice of all cinema halls and nil as error if successful
	return allCinemaHalls, nil
}

// RetrieveCinemaHallInfoByID retrieves cinema hall information by its unique ID for the admin page.
//
// Parameters:
//   - cinemaHallID (int): The unique ID of the cinema hall whose information is to be retrieved.
//
// Returns:
//   - CinemaHallForAdmin: The information of the cinema hall corresponding to the given ID.
//   - error: Returns an error if there's an issue retrieving or scanning the data, or if no cinema hall is found.
func (psql *Postgres) RetrieveCinemaHallInfoByID(cinemaHallID int) (CinemaHallForAdmin, error) {
	// SQL query to retrieve cinema hall information by ID from the database
	stmt := `SELECT * FROM cinema_hall WHERE cinema_hall_id = $1`

	// Variable to hold the cinema hall data
	var cinemaHall CinemaHallForAdmin

	// Execute the query and scan the result into the cinemaHall struct
	err := psql.DB.QueryRow(stmt, cinemaHallID).Scan(&cinemaHall.CinemaHallID, &cinemaHall.HallName, &cinemaHall.HallType, &cinemaHall.Capacity)
	if err != nil {
		// Check if no rows were found and return a custom error if so
		if errors.Is(err, sql.ErrNoRows) {
			return CinemaHallForAdmin{}, ErrCinemaHallNotFound
		}
		// Return a wrapped error if there’s an issue executing the query or scanning the data
		return CinemaHallForAdmin{}, fmt.Errorf("failed to retrieve cinema hall information for admin page: %w", err)
	}

	// Return the cinema hall information if successful
	return cinemaHall, nil
}

// DeleteCinemaHallByID deletes a cinema hall by its unique ID.
//
// Parameters:
//   - cinemaHallID (int): The unique ID of the cinema hall to be deleted.
//
// Returns:
//   - error: Returns an error if there’s an issue deleting the cinema hall or if no rows are affected.
func (psql *Postgres) DeleteCinemaHallByID(cinemaHallID int) error {
	// SQL query to delete a cinema hall by its unique ID
	stmt := `DELETE FROM cinema_hall WHERE cinema_hall_id = $1`

	// Execute the delete query
	result, err := psql.DB.Exec(stmt, cinemaHallID)
	if err != nil {
		// Return a wrapped error if query execution fails
		return fmt.Errorf("failed to delete cinema hall by ID: %w", err)
	}

	// Check the number of rows affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Return an error if checking the rows affected fails
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// If no rows were affected, return a custom error indicating the cinema hall was not found
	if rowsAffected == 0 {
		return ErrCinemaHallNotFound
	}

	// Return nil if the deletion was successful
	return nil
}

// CountCinemaSeatsByHallID retrieves the total count of cinema seats in a specific cinema hall.
//
// Parameters:
//   - hallID (int): The unique ID of the cinema hall whose seats are to be counted.
//
// Returns:
//   - int: The total count of cinema seats in the specified hall.
//   - error: Returns an error if there’s an issue with the database query.
func (psql *Postgres) CountCinemaSeatsByHallID(hallID int) (int, error) {
	// SQL query to count the number of seats in the specified cinema hall
	stmt := `SELECT COUNT(*) FROM cinema_seat WHERE hall_id = $1`

	// Execute the query and scan the result into the count variable
	var count int
	err := psql.DB.QueryRow(stmt, hallID).Scan(&count)
	if err != nil {
		// Return a wrapped error with the hallID if the query fails
		return 0, fmt.Errorf("failed to count cinema seats for hall_id %d: %w", hallID, err)
	}

	// Return the count of seats
	return count, nil
}

// InsertCinemaSeats inserts a new cinema seat into the database for a specific cinema hall.
//
// Parameters:
//   - seatRow (string): The row where the seat is located in the cinema hall.
//   - seatNumber (int): The unique number assigned to the seat within the row.
//   - seatType (string): The type of the seat (e.g., standard, VIP, etc.).
//   - hallID (int): The ID of the cinema hall where the seat is located.
//
// Returns:
//   - error: Returns an error if there’s an issue inserting the seat or if any database constraints are violated.
func (psql *Postgres) InsertCinemaSeats(seatRow string, seatNumber int, seatType string, hallID int) error {
	// SQL query to insert a new cinema seat into the cinema_seat table
	stmt := `INSERT INTO cinema_seat (seat_row, seat_number, seat_type, hall_id) VALUES ($1, $2, $3, $4)`

	// Execute the insert query
	_, err := psql.DB.Exec(stmt, seatRow, seatNumber, seatType, hallID)
	if err != nil {
		// Check for a unique violation error (23505) - seat already exists
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrCinemaSeatAlreadyExists // Custom error for seat already existing
		}

		// Check for a foreign key violation error (23503) - hall ID not found
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return ErrCinemaHallNotFound // Custom error for the hall ID not found
		}

		// Return a wrapped error if another issue occurs during the insert operation
		return fmt.Errorf("failed to insert new cinema seats: %w", err)
	}

	// Return nil if the insertion was successful
	return nil
}

// RetrieveALLCinemaSeatsByHallID retrieves all cinema seats for a given cinema hall by its ID.
//
// Parameters:
//   - hallID (int): The unique ID of the cinema hall whose seats are to be retrieved.
//
// Returns:
//   - []CinemaSeatForAdmin: A slice of CinemaSeatForAdmin containing details of all seats in the specified hall.
//   - error: Returns an error if there is an issue with the database query or no seats are found.
func (psql *Postgres) RetrieveALLCinemaSeatsByHallID(hallID int) ([]CinemaSeatForAdmin, error) {
	// SQL query to retrieve all cinema seats for a specific cinema hall
	stmt := `SELECT cinema_seat_id, seat_row, seat_number, seat_type, hall_id FROM cinema_seat WHERE hall_ID = $1`

	// Execute the query and retrieve the rows
	rows, err := psql.DB.Query(stmt, hallID)
	if err != nil {
		// Return a wrapped error if the query fails
		return nil, fmt.Errorf("failed to retrieve all cinema seats for admin page: %w", err)
	}

	// Ensure the rows are closed after processing
	defer rows.Close()

	// Slice to store all the retrieved cinema seats
	var allCinemaSeats []CinemaSeatForAdmin

	// Loop through the rows and scan each cinema seat into the struct
	for rows.Next() {
		var cinemaSeat CinemaSeatForAdmin

		err := rows.Scan(&cinemaSeat.CinemaSeatID, &cinemaSeat.SeatRow, &cinemaSeat.SeatNumber, &cinemaSeat.SeatType, &cinemaSeat.HallID)
		if err != nil {
			// If no rows are found, return an error
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrCinemaSeatNotFound
			}
			// Return a wrapped error if scanning fails
			return nil, fmt.Errorf("failed to scan cinema seats: %w", err)
		}

		// Append the scanned cinema seat to the slice
		allCinemaSeats = append(allCinemaSeats, cinemaSeat)
	}

	// If no seats were found, return an error
	if len(allCinemaSeats) == 0 {
		return nil, ErrCinemaSeatNotFound
	}

	// Return the slice of all cinema seats
	return allCinemaSeats, nil
}

// DeleteCinemaSeatByID deletes a cinema seat by its unique ID.
//
// Parameters:
//   - cinemaSeatID (int): The unique ID of the cinema seat to be deleted.
//
// Returns:
//   - error: Returns an error if there is an issue with the database query, or if no rows were affected (seat not found).
func (psql *Postgres) DeleteCinemaSeatByID(cinemaSeatID int) error {
	// SQL query to delete a cinema seat based on its unique ID
	stmt := `DELETE FROM cinema_seat WHERE cinema_seat_id = $1`

	// Execute the delete query
	result, err := psql.DB.Exec(stmt, cinemaSeatID)
	if err != nil {
		// Return a wrapped error if the query fails
		return fmt.Errorf("failed to delete cinema seat: %w", err)
	}

	// Check how many rows were affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Return a wrapped error if there’s an issue checking the affected rows
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	// If no rows were affected, the cinema seat doesn't exist in the database
	if rowsAffected == 0 {
		return ErrCinemaSeatNotFound // Custom error for when the cinema seat isn't found
	}

	// Return nil if the seat was successfully deleted
	return nil
}

// InsertNewShow inserts a new show into the database and returns the generated show ID.
//
// Parameters:
//   - showDate (string): The date of the show in string format (e.g., "2025-02-14").
//   - startTime (string): The start time of the show in string format (e.g., "14:00").
//   - hallID (int): The unique ID of the cinema hall where the show will be held.
//   - movieID (int): The unique ID of the movie being shown.
//
// Returns:
//   - showID (int): The unique ID of the newly inserted show.
//   - error: Returns an error if there is an issue with the database query or foreign key violations.
func (psql *Postgres) InsertNewShow(showDate string, startTime string, hallID int, movieID int) (int, error) {
	// SQL query to insert a new show and return the generated show_id
	stmt := `INSERT INTO show (show_date, start_time, hall_id, movie_id) VALUES ($1, $2, $3, $4) RETURNING show_id`

	var showID int

	// Execute the query to insert a new show and return the show_id
	err := psql.DB.QueryRow(stmt, showDate, startTime, hallID, movieID).Scan(&showID)
	if err != nil {
		// Handle duplicate key error (unique constraint violation)
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			// This error happens when there's a violation of the unique constraint, e.g., a show with the same hall, date, and time already exists
			return 0, ErrShowAlreadyExists
		}

		// Handle foreign key constraint violations for hall_id and movie_id
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			// Check for a missing cinema hall ID
			if pqErr.Detail != "" {
				if pqErr.Detail == "Key (hall_id)=(<hall_id>) is not present in table \"hall\"" {
					return 0, ErrCinemaHallNotFound // Custom error for hall not found
				} else if pqErr.Detail == "Key (movie_id)=(<movie_id>) is not present in table \"movie\"" {
					return 0, ErrMovieNotFoundByID // Custom error for movie not found
				}
			}

		}

		// Return a wrapped error if there is any issue with the query
		return 0, fmt.Errorf("failed to insert new show: %w", err)
	}

	// Return the generated show ID after a successful insertion
	return showID, nil
}

// RetrieveAllShowsForAdmin retrieves all shows from the database for the admin page.
//
// Returns:
//   - shows ([]ShowForAdmin): A slice of ShowForAdmin structs representing all shows.
//   - error: An error if there is any issue during the database query or row scanning.
func (psql *Postgres) RetrieveAllShowsForAdmin() ([]ShowForAdmin, error) {
	// SQL query to retrieve all shows with their show_id, show_date, start_time, hall_id, and movie_id
	stmt := `SELECT show_id, show_date, start_time, hall_id, movie_id FROM show`

	// Execute the query to get the rows
	rows, err := psql.DB.Query(stmt)
	if err != nil {
		// Return a wrapped error if the query fails
		return nil, fmt.Errorf("failed to retrieve shows: %w", err)
	}

	// Ensure the rows are closed after the function finishes
	defer rows.Close()

	// Slice to hold all shows
	var shows []ShowForAdmin

	// Iterate through the rows and scan data into the ShowForAdmin struct
	for rows.Next() {
		var show ShowForAdmin
		// Scan the row into the show struct
		if err := rows.Scan(&show.ShowID, &show.ShowDate, &show.StartTime, &show.HallID, &show.MovieID); err != nil {
			// Handle scanning errors
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrShowNotFound // Custom error if no shows are found
			}
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		// Append the successfully scanned show to the slice
		shows = append(shows, show)
	}

	// Return the slice of shows
	return shows, nil
}

// UpdateShowByID updates the details of a show in the database given a show ID.
//
// Parameters:
//   - showID (int): The ID of the show to be updated.
//   - showDate (string): The new date for the show.
//   - startTime (string): The new start time for the show.
//   - hallID (int): The ID of the cinema hall where the show will take place.
//   - movieID (int): The ID of the movie to be shown.
//
// Returns:
//   - error: If there is any issue during the update process, an error is returned.
func (psql *Postgres) UpdateShowByID(showID int, showDate string, startTime string, hallID int, movieID int) error {
	// SQL query to update the show details
	stmt := `UPDATE show SET show_date = $1, start_time = $2, hall_id = $3, movie_id = $4 WHERE show_id = $5`

	// Execute the query with the provided parameters
	result, err := psql.DB.Exec(stmt, showDate, startTime, hallID, movieID, showID)
	if err != nil {
		// Check if the error is related to foreign key constraint violations
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			if pqErr.Detail != "" {
				// If the hall ID is invalid (not present in the hall table), return a custom error
				if pqErr.Detail == fmt.Sprintf("Key (hall_id)=(%d) is not present in table \"hall\"", hallID) {
					return ErrCinemaHallNotFound
				}
				// If the movie ID is invalid (not present in the movie table), return a custom error
				if pqErr.Detail == fmt.Sprintf("Key (movie_id)=(%d) is not present in table \"movie\"", movieID) {
					return ErrMovieNotFoundByID
				}
			}
		}

		// If the error is not related to foreign key constraints, wrap and return it
		return fmt.Errorf("failed to update show with ID %d: %w", showID, err)
	}

	// Check how many rows were affected by the update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	// If no rows were affected, return an error indicating the show ID was not found
	if rowsAffected == 0 {
		return ErrShowNotFound
	}

	// Return nil if the update was successful
	return nil
}

// DeleteShowByID deletes a show from the database given the show ID.
//
// Parameters:
//   - showID (int): The ID of the show to be deleted.
//
// Returns:
//   - error: If an issue occurs while deleting the show, an error is returned.
func (psql *Postgres) DeleteShowByID(showID int) error {
	// SQL query to delete a show based on show_id
	stmt := `DELETE FROM show WHERE show_id = $1`

	// Execute the DELETE query with the provided showID
	result, err := psql.DB.Exec(stmt, showID)
	if err != nil {
		return fmt.Errorf("failed to delete show: %w", err)
	}

	// Check how many rows were affected by the delete operation
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	// If no rows were affected, return an error indicating that the show ID was not found
	if rowsAffected == 0 {
		return ErrShowNotFound
	}

	// Return nil if the show was successfully deleted
	return nil
}

// InsertNewShowSeat inserts a new show seat record into the database.
// Parameters:
//   - seatStatus (string): The status of the seat (e.g., "available", "reserved", etc.)
//   - seatPrice (int): The price of the seat for the show
//   - cinemaSeatID (int): The ID of the corresponding cinema seat
//   - showID (int): The ID of the show the seat is associated with
//
// Returns:
//   - error: If an issue occurs during the insertion, an error is returned.
func (psql *Postgres) InsertNewShowSeat(seatStatus string, seatPrice int, cinemaSeatID int, showID int) error {
	// SQL query to insert a new show seat into the show_seat table
	stmt := `INSERT INTO show_seat (cinema_seat_id, status, price, show_id) VALUES ($1, $2, $3, $4)`

	// Execute the query
	_, err := psql.DB.Exec(stmt, cinemaSeatID, seatStatus, seatPrice, showID)
	if err != nil {
		// If an error occurs during execution, return a wrapped error
		return fmt.Errorf("failed to insert new show seat: %w", err)
	}

	// Return nil if the insertion is successful
	return nil
}

// RetrieveAllShowSeats retrieves all show seats for a particular show by showID
// Parameters:
//   - showID (int): The ID of the show for which seats are being retrieved
//
// Returns:
//   - []ShowSeatForAdmin: A list of show seats associated with the show
//   - error: Returns an error if something goes wrong while fetching or processing the data
func (psql *Postgres) RetrieveAllShowSeats(showID int) ([]ShowSeatForAdmin, error) {
	// SQL query to fetch all show seats for a given show
	stmt := `SELECT show_seat_id, cinema_seat_id, status, price, show_id FROM show_seat WHERE show_id = $1`

	// Execute the query and fetch rows
	rows, err := psql.DB.Query(stmt, showID)
	if err != nil {
		// Return an error if the query fails
		return nil, fmt.Errorf("failed to retrieve all show seats for admin page: %w", err)
	}
	defer rows.Close() // Ensure that the rows are closed once the function exits

	var allShowSeats []ShowSeatForAdmin

	// Loop through all rows in the result set
	for rows.Next() {
		var showSeat ShowSeatForAdmin

		// Scan each row into the showSeat struct
		err := rows.Scan(&showSeat.ShowSeatID, &showSeat.CinemaSeatID, &showSeat.SeatStatus, &showSeat.SeatPrice, &showSeat.ShowID)
		if err != nil {
			// If scanning fails, return an error
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrShowSeatNotFound
			}
			return nil, fmt.Errorf("failed to scan show seat: %w", err)
		}

		// Append the scanned show seat to the result slice
		allShowSeats = append(allShowSeats, showSeat)
	}

	// Return the list of all show seats for the specified show
	return allShowSeats, nil
}

// UpdateShowSeatByID updates the price of a specific show seat by its ID.
// Parameters:
//   - seatPrice (int): The new price for the seat
//   - showSeatID (int): The ID of the show seat to update
//
// Returns:
//   - error: If any error occurs during the update or if no rows are affected.
func (psql *Postgres) UpdateShowSeatByID(seatPrice int, showSeatID int) error {
	// SQL query to update the price of a show seat by its ID
	stmt := `UPDATE show_seat SET price = $1 WHERE show_seat_id = $2`

	// Execute the query
	result, err := psql.DB.Exec(stmt, seatPrice, showSeatID)
	if err != nil {
		// Return error if query execution fails
		return fmt.Errorf("error occurred while updating seat price: %w", err)
	}

	// Check if any rows were affected by the update
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		// Return error if we can't check the number of affected rows
		return fmt.Errorf("error occurred while checking affected rows: %w", err)
	}

	// If no rows were affected, return a custom error indicating the seat wasn't found
	if rowsAffected == 0 {
		return ErrShowSeatNotFound
	}

	// Return nil if the update was successful
	return nil
}
