package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type DBContractMovies interface {
	RetrieveAllCarouselImages() ([]CarouselImage, error)
	RetrieveAllMovies() ([]AllMovies, error)
	RetrieveAMovie(movieID int) (Movie, error)
	RetrieveAllActorsCrewsByMovieID(movieID int) ([]ActorsCrewsOfMovie, error)
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

// RetrieveAllMovies retrieves all movie records from the database.
//
// Parameters:
//   - psql: a pointer to the Postgres struct, which contains the database connection
//     and is responsible for querying the database.
//
// Returns:
//   - A slice of AllMovies structs containing movie details like ID, title, language,
//     poster URL, rating, rating provider, and age limit.
//   - An error if the database query fails or if there are issues scanning the results.
func (psql *Postgres) RetrieveAllMovies() ([]AllMovies, error) {
	stmt := `SELECT id, title, genre, language, poster_url, rating, rating_provider, age_limit FROM movies`

	// Execute the SQL query to fetch all movies from the database
	rows, err := psql.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all movies from the database: %w", err)
	}

	defer rows.Close()

	var movies []AllMovies

	// Loop through the rows and scan the values into AllMovies struct
	for rows.Next() {
		var movie AllMovies
		err := rows.Scan(&movie.ID, &movie.Title, &movie.Genre, &movie.Language, &movie.PosterUrl, &movie.Rating, &movie.RatingProvider, &movie.AgeLimit)
		if err != nil {
			return nil, fmt.Errorf("failed to scan all movies: %w", err)
		}
		movies = append(movies, movie)
	}

	// Return the slice of movies
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
func (psql *Postgres) RetrieveAMovie(movieID int) (Movie, error) {
	// SQL query to fetch a movie by its ID
	stmt := `SELECT id, title, description, genre, language, trailer_url, poster_url, rating, rating_provider, duration, release_date, age_limit FROM movies WHERE id = $1`

	// Execute the query and get the result
	row := psql.DB.QueryRow(stmt, movieID)

	var movie Movie

	// Scan the row into the movie struct
	err := row.Scan(&movie.ID, &movie.Title, &movie.Description, &movie.Genre, &movie.Language, &movie.TrailerUrl, &movie.PosterUrl, &movie.Rating, &movie.RatingProvider, &movie.Duration, &movie.ReleaseDate, &movie.AgeLimit)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Return a specific error if no movie is found
			return Movie{}, ErrMovieNotFoundByID
		} else {
			// Return a generic error if there are any other issues
			return Movie{}, fmt.Errorf("failed to retrieve a movie by id from the database: %w", err)
		}
	}

	// Return the movie if found
	return movie, nil
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
