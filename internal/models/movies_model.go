package models

import "fmt"

type DBContractMovies interface {
	RetrieveAllCarouselImages() ([]CarouselImage, error)
	RetrieveAllMovies() ([]AllMovies, error)
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
// - psql: a pointer to the Postgres struct, which contains the database connection
//         and is responsible for querying the database.
//
// Returns:
// - A slice of AllMovies structs containing movie details like ID, title, language,
//   poster URL, rating, rating provider, and age limit.
// - An error if the database query fails or if there are issues scanning the results.
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
