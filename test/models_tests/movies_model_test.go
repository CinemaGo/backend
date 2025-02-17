package modelstests

import (
	"cinemaGo/backend/internal/models"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestRetrieveAllCarouselImages(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	defer db.Close()

	psql := &models.Postgres{DB: db}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "image_url"}).AddRow(1, "https://example.com/image1.jpg").AddRow(2, "https://example.com/image2.jpg")
		mock.ExpectQuery("SELECT id, image_url FROM carousel_images").WillReturnRows(rows)

		carouselImages, err := psql.RetrieveAllCarouselImages()

		assert.NoError(t, err)
		assert.Len(t, carouselImages, 2)
		assert.Equal(t, carouselImages[0].ID, 1)
		assert.Equal(t, carouselImages[0].ImageURL, "https://example.com/image1.jpg")
		assert.Equal(t, carouselImages[1].ID, 2)
		assert.Equal(t, carouselImages[1].ImageURL, "https://example.com/image2.jpg")
	})

	t.Run("query_error", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, image_url FROM carousel_images").WillReturnError(fmt.Errorf("query failed"))

		carouselImages, err := psql.RetrieveAllCarouselImages()

		assert.Error(t, err)
		assert.Nil(t, carouselImages)
		assert.Contains(t, err.Error(), "failed to retrieve all carousel images")
	})

	t.Run("no_result", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, image_url FROM carousel_images").WillReturnRows(sqlmock.NewRows([]string{"id", "image_url"}))

		carouselImages, err := psql.RetrieveAllCarouselImages()

		assert.NoError(t, err)
		assert.Empty(t, carouselImages)
	})

	t.Run("scan_error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "image_url"}).AddRow("invalid", "https://example.com/image1.jpg")
		mock.ExpectQuery("SELECT id, image_url FROM carousel_images").WillReturnRows(rows)

		carouselImages, err := psql.RetrieveAllCarouselImages()

		assert.Error(t, err)
		assert.Nil(t, carouselImages)
		assert.Contains(t, err.Error(), "failed to scan all carousel images")
	})
}

func TestRetrieveAllShowsMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock database: %v", err)
	}

	defer db.Close()

	psql := &models.Postgres{DB: db}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"show_id", "movie_id", "movie_title", "movie_genre", "movie_language", "movie_poster_url", "movie_rating", "movie_rating_provider", "movie_age_limit"}).
			AddRow(1, 101, "Movie 1", "Action", "English", "https://example.com/poster1.jpg", 8.5, "IMDB", "UA").
			AddRow(2, 102, "Movie 2", "Comedy", "Spanish", "https://example.com/poster2.jpg", 7.0, "Rotten Tomatoes", "A")

		mock.ExpectQuery("SELECT s.show_id AS show_id, m.id AS movie_id").WillReturnRows(rows)

		movies, err := psql.RetrieveAllShowsMovie()

		assert.NoError(t, err)
		assert.Len(t, movies, 2)
		assert.Equal(t, movies[0].ShowID, 1)
		assert.Equal(t, movies[0].MovieID, 101)
		assert.Equal(t, movies[0].MovieTitle, "Movie 1")
		assert.Equal(t, movies[0].MovieGenre, "Action")
		assert.Equal(t, movies[0].MovieLanguage, "English")
		assert.Equal(t, movies[0].MoviePosterUrl, "https://example.com/poster1.jpg")
		assert.Equal(t, movies[0].MovieRating, float32(8.5))
		assert.Equal(t, movies[0].MovieRatingProvider, "IMDB")
		assert.Equal(t, movies[0].MovieAgeLimit, "UA")
	})

	t.Run("query_error", func(t *testing.T) {
		mock.ExpectQuery("SELECT s.show_id AS show_id, m.id AS movie_id").WillReturnError(fmt.Errorf("query failed"))

		movies, err := psql.RetrieveAllShowsMovie()

		assert.Error(t, err)
		assert.Nil(t, movies)
		assert.Contains(t, err.Error(), "failed to retrieve all movies")
	})

	t.Run("no_result", func(t *testing.T) {
		mock.ExpectQuery("SELECT s.show_id AS show_id, m.id AS movie_id").WillReturnRows(sqlmock.NewRows([]string{"show_id", "movie_id", "movie_title", "movie_genre", "movie_language", "movie_poster_url", "movie_rating", "movie_rating_provider", "movie_age_limit"}))

		movies, err := psql.RetrieveAllShowsMovie()

		assert.NoError(t, err)
		assert.Empty(t, movies)
	})

	t.Run("scan_error", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"show_id", "movie_id", "movie_title", "movie_genre", "movie_language", "movie_poster_url", "movie_rating", "movie_rating_provider", "movie_age_limit"}).
			AddRow(1, "invalid_id", "Movie 1", "Action", "English", "https://example.com/poster1.jpg", 8.5, "IMDB", "UA")

		mock.ExpectQuery("SELECT s.show_id AS show_id, m.id AS movie_id").WillReturnRows(rows)

		movies, err := psql.RetrieveAllShowsMovie()

		assert.Error(t, err)
		assert.Nil(t, movies)
		assert.Contains(t, err.Error(), "failed to scan all movies")
	})
}
