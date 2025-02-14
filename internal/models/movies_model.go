package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type DBContractMovies interface {
	RetrieveAllCarouselImages() ([]CarouselImage, error)
	RetrieveAllShowsMovie() ([]AllShowsMovie, error)
	RetrieveAShowMovie(showID int) (AShowMovie, error)
	RetrieveAllActorsCrewsByMovieID(movieID int) ([]ActorsCrewsOfMovie, error)
	RetriveActorCrewInfo(actorCrewID int) (ActorCrewInfo, error)
	RetrieveMoviesByActorCrewID(actorCrewID int) ([]ActorCrewMovies, error)
}

type Movies struct {
	db DBContractMovies
}

func NewMovies(db DBContractMovies) *Movies {
	return &Movies{db: db}
}

// RetrieveAllCarouselImages retrieves all carusel images from the database:
//
// Parametrs:
// - psql: apointer to the Postgres struct that contains the database connection.
//
// Returns:
// - A slice of CarouselImage structs, each containing the ID and ImageURL of a carousel image.
// - An error if the database query fails or if there are issues scanning the data.
func (psql *Postgres) RetrieveAllCarouselImages() ([]CarouselImage, error) {
	stmt := `SELECT id, image_url FROM carousel_images`

	// Execute the query to retrieve the carousel images from the database
	rows, err := psql.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all carousel images from the database: %w", err)
	}

	defer rows.Close()

	var carouselImages []CarouselImage

	// Iterate through the rows returned by the query
	for rows.Next() {
		var carouselImage CarouselImage

		// Scan the columns into the CarouselImage struct
		err := rows.Scan(&carouselImage.ID, &carouselImage.ImageURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan all carousel images: %w", err)
		}
		carouselImages = append(carouselImages, carouselImage)
	}
	// Return the slice of carousel images
	return carouselImages, nil
}
// RetrieveAllShowsMovie retrieves all shows with their associated movie details from the database.
//
// It performs a SQL query to fetch information from the 'show' table and the 'movies' table,
// joining them on the movie ID to retrieve details about each movie along with its show information.
// 
// Returns:
//   - []AllShowsMovie: A slice of `AllShowsMovie` structs containing the show and movie details.
//   - error: If any error occurs during the execution of the query or scanning of rows, it returns an error.
func (psql *Postgres) RetrieveAllShowsMovie() ([]AllShowsMovie, error) {
	// SQL query that joins the 'show' and 'movies' tables to retrieve show and movie details
	stmt := `SELECT s.show_id AS show_id, m.id AS movie_id, m.title AS movie_title, m.genre AS movie_genre, m.language AS movie_language, m.poster_url AS movie_poster_url, m.rating AS movie_rating, m.rating_provider AS movie_rating_provider, m.age_limit AS movie_age_limit FROM show s JOIN movies m ON s.movie_id = m.id`

	// Execute the query
	rows, err := psql.DB.Query(stmt)
	if err != nil {
		// If there's an error executing the query, return a wrapped error
		return nil, fmt.Errorf("failed to retrieve all movies from the database: %w", err)
	}

	// Ensure rows are properly closed after function execution
	defer rows.Close()

	// Slice to hold the retrieved movies
	var movies []AllShowsMovie

	// Iterate over the rows returned from the query
	for rows.Next() {
		// Variable to hold the data of a single movie
		var movie AllShowsMovie

		// Scan the row data into the 'movie' variable
		err := rows.Scan(&movie.ShowID, &movie.MovieID, &movie.MovieTitle, &movie.MovieGenre, &movie.MovieLanguage, &movie.MoviePosterUrl, &movie.MovieRating, &movie.MovieRatingProvider, &movie.MovieAgeLimit)
		if err != nil {
			// If scanning fails, return an error with a wrapped message
			return nil, fmt.Errorf("failed to scan all movies: %w", err)
		}

		// Append the scanned movie details to the 'movies' slice
		movies = append(movies, movie)
	}

	// Return the list of movies after processing all rows
	return movies, nil
}


// RetrieveAMovie retrieves a movie by its ID from the database.
//
// Parameters:
// - id: The unique identifier for the movie to be retrieved.
//
// Returns:
// - A Movie struct containing all the movie details if found.
// - An error if the movie is not found (ErrMovieNotFoundByID) or if there's an issue querying the database.
func (psql *Postgres) RetrieveAShowMovie(showID int) (AShowMovie, error) {
	// SQL query to fetch a movie by its ID
	stmt := `SELECT s.show_id AS show_id, m.id AS movie_id, m.title AS movie_title, m.description AS movie_description, m.genre AS movie_genre, m.language AS movie_language, m.trailer_url AS movie_trailer_url, m.poster_url AS movie_poster_url, m.rating AS movie_rating, m.rating_provider AS movie_rating_provider, m.duration AS movie_duration, m.release_date AS movie_release_date, m.age_limit AS movie_age_limit FROM show s JOIN movies m ON s.movie_id = m.id WHERE s.show_id = $1`

	// Execute the query and get the result
	row := psql.DB.QueryRow(stmt, showID)

	var aShowMovie AShowMovie

	// Scan the row into the movie struct
	err := row.Scan(&aShowMovie.ShowID,&aShowMovie.MovieID, &aShowMovie.MovieTitle, &aShowMovie.MovieDescription, &aShowMovie.MovieGenre, &aShowMovie.MovieLanguage, &aShowMovie.MovieTrailerUrl, &aShowMovie.MoviePosterUrl, &aShowMovie.MovieRating, &aShowMovie.MovieRatingProvider, &aShowMovie.MovieDuration, &aShowMovie.MovieReleaseDate, &aShowMovie.MovieAgeLimit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return a specific error if no movie is found
			return AShowMovie{}, ErrMovieNotFoundByID
		} else {
			// Return a generic error if there are any other issues
			return AShowMovie{}, fmt.Errorf("failed to retrieve a movie by id from the database: %w", err)
		}
	}

	// Return the movie if found
	return aShowMovie, nil
}

// RetrieveAllActorsCrewsByMovieID retrieves all actors and crew members for a specific movie
// by querying the database based on the provided movie ID.
//
// Parameters:
// - movieID: The ID of the movie whose actors and crew are being retrieved.
//
// Returns:
// - A slice of ActorsCrewsOfMovie containing the details of each actor/crew member.
// - An error if there are issues retrieving or scanning the data.
func (psql *Postgres) RetrieveAllActorsCrewsByMovieID(movieID int) ([]ActorsCrewsOfMovie, error) {
	// SQL query to retrieve actors and crew members for a given movie ID
	stmt := `SELECT ac.id, ac.full_name, ac.image_url, ac.role_description, ac.is_actor FROM actors_crew ac JOIN movie_actors_crew mac ON ac.id = mac.actor_crew_id JOIN movies m ON mac.movie_id = m.id WHERE m.id = $1`

	// Execute the query
	rows, err := psql.DB.Query(stmt, movieID)
	if err != nil {
		// Handle any errors executing the query
		return nil, fmt.Errorf("failed to retrieve all actors and crews from database: %w", err)
	}
	defer rows.Close() // Ensure rows are closed after function returns

	// Slice to hold all the actors/crew details
	var allActorsCrew []ActorsCrewsOfMovie

	// Iterate through the result set
	for rows.Next() {
		var actorCrew ActorsCrewsOfMovie
		// Scan the row into the actorCrew struct
		err := rows.Scan(&actorCrew.ID, &actorCrew.FullName, &actorCrew.ImageURL, &actorCrew.RoleDescription, &actorCrew.IsActor)
		if err != nil {
			// If no rows are found for the given movie ID, return an appropriate error
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrMovieNotFoundByID
			}
			// Handle any errors scanning the result rows
			return nil, fmt.Errorf("failed to scan all actors and crew: %w", err)
		}

		// Append the actor/crew member to the result slice
		allActorsCrew = append(allActorsCrew, actorCrew)
	}

	// Return the result slice
	return allActorsCrew, nil
}

// RetriveActorCrewInfo retrieves detailed information about an actor or crew member
// from the database based on the provided actorCrewID.
//
// Parameters:
// - actorCrewID: The ID of the actor or crew member whose details are being retrieved.
//
// Returns:
// - ActorCrewInfo: A struct containing detailed information about the actor/crew member.
// - An error if there are issues with the database query or scanning the data.
func (psql *Postgres) RetriveActorCrewInfo(actorCrewID int) (ActorCrewInfo, error) {
	// SQL query to retrieve detailed actor/crew info by ID
	stmt := `SELECT id, full_name, image_url, occupation, role_description, born_date, birthplace, about FROM actors_crew WHERE id = $1`

	// Declare a variable to hold the actor/crew info
	var actorCrewInfo ActorCrewInfo
	row := psql.DB.QueryRow(stmt, actorCrewID)

	// Scan the query result into the actorCrewInfo struct
	err := row.Scan(&actorCrewInfo.ID, &actorCrewInfo.FullName, &actorCrewInfo.ImageURL, &actorCrewInfo.Occupation,
		&actorCrewInfo.RoleDescription, &actorCrewInfo.BornDate, &actorCrewInfo.Birthplace, &actorCrewInfo.About)

	// Error handling: check if no rows are returned or other scanning errors
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ActorCrewInfo{}, ErrActorCrewNotFoundByID // Return a custom error if no rows are found
		}
		return ActorCrewInfo{}, fmt.Errorf("failed to scan actor or crew info: %w", err) // Wrap other errors with context
	}

	// Return the actor/crew member's information if successful
	return actorCrewInfo, nil
}

// RetrieveMoviesByActorCrewID retrieves all movies associated with a specific actor or crew member,
// based on the provided actorCrewID.
//
// Parameters:
// - actorCrewID (int): The ID of the actor or crew member whose associated movies are being fetched.
//
// Returns:
// - []ActorCrewMovies: A slice of ActorCrewMovies structs containing movie details for the given actor or crew member.
// - error: Returns nil if the retrieval was successful, or an error if any occurred during the process.
func (psql *Postgres) RetrieveMoviesByActorCrewID(actorCrewID int) ([]ActorCrewMovies, error) {
	// SQL query to fetch all movies associated with a specific actor or crew member
	stmt := `SELECT m.id AS movie_id, m.title AS movie_title, m.poster_url FROM movies m JOIN movie_actors_crew mac ON m.id = mac.movie_id JOIN actors_crew ac ON mac.actor_crew_id = ac.id WHERE ac.id = $1`

	// Execute the query and handle any errors
	rows, err := psql.DB.Query(stmt, actorCrewID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all movies of actors or crews: %w", err)
	}

	// Ensure the rows are closed when done
	defer rows.Close()

	// Declare a slice to store the actor/crew's movies
	var actorCrewMovies []ActorCrewMovies

	// Iterate over the rows and scan the data into the struct
	for rows.Next() {
		var acMovie ActorCrewMovies

		// Scan the row data into the struct fields
		err := rows.Scan(&acMovie.ID, &acMovie.Title, &acMovie.PosterUrl)
		if err != nil {
			// If no rows were found for the actor/crew ID, return a specific error
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrActorCrewNotFoundByID
			}
			// Wrap any other errors with context for better debugging
			return nil, fmt.Errorf("failed to scan a movie of the actor or crew: %w", err)
		}

		// Append the movie data to the slice
		actorCrewMovies = append(actorCrewMovies, acMovie)
	}

	// Return the list of actor/crew movies and nil for the error if everything was successful
	return actorCrewMovies, nil
}
