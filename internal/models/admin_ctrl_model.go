package models

import (
	"database/sql"
	"errors"
	"fmt"

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
	UpdateCinemaHallInfoByID(hallName, hallType string, capacity, cinemaHallID int) error
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

func (psql *Postgres) InsertNewCarouselImages(imageURL, title, description string, orderPriority int) error {
	stmt := `INSERT INTO carousel_images (image_url, title, description, order_priority) VALUES ($1, $2, $3, $4)`

	_, err := psql.DB.Exec(stmt, imageURL, title, description, orderPriority)
	if err != nil {
		return fmt.Errorf("failed to insert data into carousel images: %w", err)
	}
	return nil
}

func (psql *Postgres) RetrieveCarouselImagesForAdmin() ([]CarouselImageForAdmin, error) {
	stmt := `SELECT *  FROM carousel_images`

	rows, err := psql.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve carousel images for admin page: %w", err)
	}

	defer rows.Close()

	var carouselImages []CarouselImageForAdmin

	for rows.Next() {
		var carouselImage CarouselImageForAdmin

		err := rows.Scan(&carouselImage.CarouselImageID, &carouselImage.ImageURL, &carouselImage.Title, &carouselImage.Description, &carouselImage.OrderPriority, &carouselImage.CreatedAt, &carouselImage.UpdatedAt)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrAdminPageCarouselImagesNotFound
			}
			return nil, fmt.Errorf("failed to retrieve admin page carousel images: %w", err)
		}

		carouselImages = append(carouselImages, carouselImage)
	}

	return carouselImages, nil
}

func (psql *Postgres) UpdateCarouselImagesByID(imageURL, title, description string, orderPriority int, carouselImageID int) error {
	stmt := `UPDATE carousel_images SET image_url = $1, title = $2, description = $3, order_priority = $4 WHERE id = $5`

	result, err := psql.DB.Exec(stmt, imageURL, title, description, orderPriority, carouselImageID)
	if err != nil {
		return fmt.Errorf("failed to insert data into carousel images: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAdminPageCarouselImagesNotFound
	}

	return nil
}

func (psql *Postgres) DeleteCarouselImagesByID(carouselImageID int) error {
	stmt := `DELETE FROM carousel_images WHERE id = $1`

	result, err := psql.DB.Exec(stmt, carouselImageID)
	if err != nil {
		return fmt.Errorf("failed to delete carousel image by id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAdminPageCarouselImagesNotFound
	}

	return nil
}

func (psql *Postgres) InsertNewMovie(title string, description, genre, language, trailerURL, posterURL string, rating int, ratingProvider string, duration int, releaseDate, ageLimit string) error {
	stmt := `INSERT INTO movies (title, description, genre, language, trailer_url, poster_url, rating, rating_provider, duration, release_date, age_limit) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := psql.DB.Exec(stmt, title, description, genre, language, trailerURL, posterURL, rating, ratingProvider, duration, releaseDate, ageLimit)
	if err != nil {
		return fmt.Errorf("failed to insert data into movies: %w", err)
	}

	return nil
}

func (psql *Postgres) RetrieveAllMoviesForAdmin() ([]AllMoviesForAdmin, error) {
	stmt := `SELECT id, title FROM movies`

	rows, err := psql.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to retriece all movies for admin page: %w", err)
	}

	defer rows.Close()

	var allMovies []AllMoviesForAdmin

	for rows.Next() {
		var movie AllMoviesForAdmin

		err := rows.Scan(&movie.ID, &movie.Title)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrAdminPageMovieNotFound
			}
			return nil, fmt.Errorf("failed to retrieve all movies for admin gape: %w", err)
		}

		allMovies = append(allMovies, movie)
	}

	return allMovies, nil
}

func (psql *Postgres) RetrieveAMovieForAdmin(movieID int) (MovieForAdmin, error) {
	stmt := `SELECT * FROM movies WHERE id = $1`

	var movie MovieForAdmin
	err := psql.DB.QueryRow(stmt, movieID).Scan(&movie.MovieID, &movie.Title, &movie.Description, &movie.Genere, &movie.Language, &movie.TrailerURL, &movie.PosterURL, &movie.Rating, &movie.RatingProvider, &movie.Duration, &movie.RelaseDate, &movie.AgeLimit, &movie.CreatedAt, &movie.UpdatesAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return MovieForAdmin{}, ErrAdminPageMovieNotFound
		}
		return MovieForAdmin{}, fmt.Errorf("failed to retrieve a movie for admin page: %w", err)
	}

	return movie, nil
}

func (psql *Postgres) UpdateMovieInfoForAdminByMovieID(movieID int, title, description, genre, language, trailerURL, posterURL string, rating float32, ratingProvider string, duration int, relaseDate string, ageLimit string) error {
	stmt := `UPDATE movies SET title = $1, description = $2, genre = $3, language = $4, trailer_url = $5, poster_url = $6, rating = $7, rating_provider = $8, duration = $9, release_date = $10, age_limit = $11, updated_at = CURRENT_TIMESTAMP WHERE id = $12;`

	result, err := psql.DB.Exec(stmt, title, description, genre, language, trailerURL, posterURL, rating, ratingProvider, duration, relaseDate, ageLimit, movieID)
	if err != nil {
		return fmt.Errorf("failed to update movie informations: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAdminPageMovieNotFound
	}

	return nil
}

func (psql *Postgres) DeleteMovieByMovieID(movieID int) error {
	stmt := `DELETE FROM movies WHERE id = $1;`

	result, err := psql.DB.Exec(stmt, movieID)
	if err != nil {
		return fmt.Errorf("failed to delete movie: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrAdminPageMovieNotFound
	}

	return nil
}

func (psql *Postgres) InsertActorsCrew(fullName, imageURL, occupation, roleDescription, bornDate, birthplace, about string, isActor bool) (int, error) {
	stmt := `INSERT INTO actors_crew (full_name, image_url, occupation, role_description, born_date, birthplace, about, is_actor) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	var actorCrewID int

	err := psql.DB.QueryRow(stmt, fullName, imageURL, occupation, roleDescription, bornDate, birthplace, about, isActor).Scan(&actorCrewID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert actors_crew info into database: %w", err)
	}

	return actorCrewID, nil
}

func (psql *Postgres) InsertMovieActorCrew(movieID, actorCrewID int) error {
	stmt := `INSERT INTO movie_actors_crew (movie_id, actor_crew_id) VALUES ($1, $2)`

	_, err := psql.DB.Exec(stmt, movieID, actorCrewID)
	if err != nil {
		return fmt.Errorf("failed to insert into movie_actors_crew: %w", err)
	}

	return nil
}

func (psql *Postgres) RetrieveAllActorsCrewByMovieID(movieID int) ([]AllActorCrewForAdmin, error) {
	stmt := `SELECT ac.id, ac.full_name, ac.image_url, ac.role_description, ac.is_actor FROM actors_crew ac JOIN movie_actors_crew mac ON ac.id = mac.actor_crew_id JOIN movies m ON mac.movie_id = m.id WHERE m.id = $1`

	rows, err := psql.DB.Query(stmt, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all actorscrew by movieID for admin page: %w", err)
	}

	defer rows.Close()

	var allActorCrew []AllActorCrewForAdmin

	for rows.Next() {
		var actorCrew AllActorCrewForAdmin

		err := rows.Scan(&actorCrew.ActorCrewID, &actorCrew.FullName, &actorCrew.ImageURL, &actorCrew.RoleDescription, &actorCrew.IsActor)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrActorCrewNotFound
			}
			return nil, fmt.Errorf("failed to scan actorscrew by movieID for admin page: %w", err)
		}

		allActorCrew = append(allActorCrew, actorCrew)
	}

	return allActorCrew, nil
}

func (psql *Postgres) RetrieveActorCrew(actorCrewID int) (ActorCrewForAdmin, error) {

	stmt := `SELECT * FROM actors_crew WHERE id = $1`
	var actorCrew ActorCrewForAdmin
	err := psql.DB.QueryRow(stmt, actorCrewID).Scan(&actorCrew.ActorCrewID, &actorCrew.FullName, &actorCrew.ImageURL, &actorCrew.Occupation, &actorCrew.RoleDescription, &actorCrew.BornDate, &actorCrew.Birthplace, &actorCrew.About, &actorCrew.IsActor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ActorCrewForAdmin{}, ErrActorCrewNotFound
		}
		return ActorCrewForAdmin{}, fmt.Errorf("failed to scan actorCrew for admin page: %w", err)
	}

	return actorCrew, nil
}

func (psql *Postgres) UpdateActorCrewInformation(fullName, imageURL, occupation, roleDescription, bornDate, birthplace, about string, isActor bool, actorCrewID int) error {

	stmt := `UPDATE actors_crew SET full_name = $1, image_url = $2, occupation = $3, role_description = $4, born_date = $5, birthplace = $6, about = $7, is_actor = $8 WHERE id = $9`

	result, err := psql.DB.Exec(stmt, fullName, imageURL, occupation, roleDescription, bornDate, birthplace, about, isActor, actorCrewID)
	if err != nil {
		return fmt.Errorf("failed to update actorCrew informations: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrActorCrewNotFound
	}

	return nil
}

func (psql *Postgres) DeleteActorCrewByID(actorCrewID int) error {
	stmt := `DELETE FROM actors_crew WHERE id = $1;`

	result, err := psql.DB.Exec(stmt, actorCrewID)
	if err != nil {
		return fmt.Errorf("failed to delete actorsCrew: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrActorCrewNotFound
	}

	return nil
}

func (psql *Postgres) InsertNewCinemaHall(hallName, hallType string, capacity int) error {
	stmt := `INSERT INTO cinema_hall (hall_name, hall_type, capacity) VALUES ($1, $2, $3)`

	_, err := psql.DB.Exec(stmt, hallName, hallType, capacity)
	if err != nil {
		return fmt.Errorf("failed to insert new hall informations: %w", err)
	}

	return nil
}

func (psql *Postgres) RetrieveAllCinemaHallsForAdmin() ([]CinemaHallForAdmin, error) {
	stmt := `SELECT * FROM cinema_hall`

	rows, err := psql.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all hall informations for admin page: %w", err)
	}

	defer rows.Close()

	var allCinemaHalls []CinemaHallForAdmin

	for rows.Next() {
		var cinemaHall CinemaHallForAdmin

		err := rows.Scan(&cinemaHall.CinemaHallID, &cinemaHall.HallName, &cinemaHall.HallType, &cinemaHall.Capacity)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrCinemaHallNotFound
			}
			return nil, fmt.Errorf("failed to scan cinema hall informations: %w", err)
		}

		allCinemaHalls = append(allCinemaHalls, cinemaHall)
	}

	return allCinemaHalls, nil
}

func (psql *Postgres) RetrieveCinemaHallInfoByID(cinemaHallID int) (CinemaHallForAdmin, error) {
	stmt := `SELECT * FROM cinema_hall WHERE cinema_hall_id = $1`

	var cinemaHall CinemaHallForAdmin

	err := psql.DB.QueryRow(stmt, cinemaHallID).Scan(&cinemaHall.CinemaHallID, &cinemaHall.HallName, &cinemaHall.HallType, &cinemaHall.Capacity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CinemaHallForAdmin{}, ErrCinemaHallNotFound
		}
		return CinemaHallForAdmin{}, fmt.Errorf("failed to retrieve a cinema hall information for admin page: %w", err)
	}

	return cinemaHall, nil
}

func (psql *Postgres) UpdateCinemaHallInfoByID(hallName, hallType string, capacity, cinemaHallID int) error {
	stmt := `UPDATE cinema_hall SET hall_name = $1, hall_type = $2, capacity = $3 WHERE cinema_hall_id = $4`

	result, err := psql.DB.Exec(stmt, hallName, hallType, capacity, cinemaHallID)
	if err != nil {
		return fmt.Errorf("failed to update cinema hall informations: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCinemaHallNotFound
	}

	return nil
}

func (psql *Postgres) DeleteCinemaHallByID(cinemaHallID int) error {
	stmt := `DELETE FROM cinema_hall WHERE cinema_hall_id = $1`

	result, err := psql.DB.Exec(stmt, cinemaHallID)
	if err != nil {
		return fmt.Errorf("failed to delete cinema hall by id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCinemaHallNotFound
	}

	return nil
}

func (psql *Postgres) CountCinemaSeatsByHallID(hallID int) (int, error) {
	stmt := `SELECT COUNT(*) FROM cinema_seat WHERE hall_id = $1`

	var count int
	err := psql.DB.QueryRow(stmt, hallID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count cinema seats for hall_id %d: %w", hallID, err)
	}

	return count, nil
}

func (psql *Postgres) InsertCinemaSeats(seatRow string, seatNumber int, seatType string, hallID int) error {
	stmt := `INSERT INTO cinema_seat (seat_row, seat_number, seat_type, hall_id) VALUES ($1, $2, $3, $4)`

	_, err := psql.DB.Exec(stmt, seatRow, seatNumber, seatType, hallID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return ErrCinemaSeatAlreadyExists
		}

		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			return ErrCinemaHallNotFound
		}

		return fmt.Errorf("failed to insert new cinema seats: %w", err)
	}

	return nil
}

func (psql *Postgres) RetrieveALLCinemaSeatsByHallID(hallID int) ([]CinemaSeatForAdmin, error) {
	stmt := `SELECT cinema_seat_id, seat_row, seat_number, seat_type, hall_id  FROM cinema_seat WHERE hall_ID = $1`

	rows, err := psql.DB.Query(stmt, hallID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all cinema seats for admin page: %w", err)
	}

	defer rows.Close()

	var allCinemaSeats []CinemaSeatForAdmin

	for rows.Next() {
		var cinemaSeat CinemaSeatForAdmin

		err := rows.Scan(&cinemaSeat.CinemaSeatID, &cinemaSeat.SeatRow, &cinemaSeat.SeatNumber, &cinemaSeat.SeatType, &cinemaSeat.HallID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrCinemaSeatNotFound
			}
			return nil, fmt.Errorf("failed to scan cinema seats: %w", err)
		}

		allCinemaSeats = append(allCinemaSeats, cinemaSeat)
	}

	if len(allCinemaSeats) == 0 {
		return nil, ErrCinemaSeatNotFound
	}

	return allCinemaSeats, nil
}

func (psql *Postgres) DeleteCinemaSeatByID(cinemaSeatID int) error {
	stmt := `DELETE FROM cinema_seat WHERE cinema_seat_id = $1`

	result, err := psql.DB.Exec(stmt, cinemaSeatID)
	if err != nil {
		return fmt.Errorf("failed to delete cinema seat: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return ErrCinemaSeatNotFound
	}

	return nil
}

func (psql *Postgres) InsertNewShow(showDate string, startTime string, hallID int, movieID int) (int, error) {
	stmt := `INSERT INTO show (show_date, start_time, hall_id, movie_id) VALUES ($1, $2, $3, $4) RETURNING show_id`

	var showID int

	err := psql.DB.QueryRow(stmt, showDate, startTime, hallID, movieID).Scan(&showID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			if pqErr.Detail != "" {
				if pqErr.Detail == "Key (hall_id)=(<hall_id>) is not present in table \"hall\"" {
					return 0, ErrCinemaHallNotFound
				} else if pqErr.Detail == "Key (movie_id)=(<movie_id>) is not present in table \"movie\"" {
					return 0, ErrMovieNotFoundByID
				}
			}
		}
		return 0, fmt.Errorf("failed to insert new show: %w", err)
	}

	return showID, nil
}

func (psql *Postgres) RetrieveAllShowsForAdmin() ([]ShowForAdmin, error) {
	stmt := `SELECT show_id, show_date, start_time, hall_id, movie_id FROM show`

	rows, err := psql.DB.Query(stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve shows: %w", err)
	}

	defer rows.Close()

	var shows []ShowForAdmin
	for rows.Next() {
		var show ShowForAdmin
		if err := rows.Scan(&show.ShowID, &show.ShowDate, &show.StartTime, &show.HallID, &show.MovieID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrShowNotFound
			}
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		shows = append(shows, show)
	}

	return shows, nil
}

func (psql *Postgres) UpdateShowByID(showID int, showDate string, startTime string, hallID int, movieID int) error {
	stmt := `UPDATE show SET show_date = $1, start_time = $2, hall_id = $3, movie_id = $4 WHERE show_id = $5`

	result, err := psql.DB.Exec(stmt, showDate, startTime, hallID, movieID, showID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23503" {
			if pqErr.Detail != "" {
				if pqErr.Detail == fmt.Sprintf("Key (hall_id)=(%d) is not present in table \"hall\"", hallID) {
					return ErrCinemaHallNotFound
					// fmt.Errorf("hall_id %d is invalid: %w", hallID, ErrInvalidHallID)
				}
				if pqErr.Detail == fmt.Sprintf("Key (movie_id)=(%d) is not present in table \"movie\"", movieID) {
					return ErrMovieNotFoundByID
					// fmt.Errorf("movie_id %d is invalid: %w", movieID, ErrInvalidMovieID)
				}
			}
		}

		return fmt.Errorf("failed to update show with ID %d: %w", showID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrShowNotFound
	}

	return nil
}

func (psql *Postgres) DeleteShowByID(showID int) error {
	stmt := `DELETE FROM show WHERE show_id = $1`

	result, err := psql.DB.Exec(stmt, showID)
	if err != nil {
		return fmt.Errorf("failed to delete show: %w", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrShowNotFound
		// fmt.Errorf("show_id %d is invalid: %w", showID, ErrShowIDNotFound)
	}

	return nil
}

func (psql *Postgres) InsertNewShowSeat(seatStatus string, seatPrice int, cinemSeatID int, showID int) error {
	stmt := `INSERT INTO show_seat (cinema_seat_id, status, price, show_id) VALUES ($1, $2, $3, $4)`

	_, err := psql.DB.Exec(stmt, cinemSeatID, seatStatus, seatPrice, showID)
	if err != nil {
		return fmt.Errorf("failed to insert new show seat: %w", err)
	}
	return nil
}

func (psql *Postgres) RetrieveAllShowSeats(showID int) ([]ShowSeatForAdmin, error) {
	stmt := `SELECT show_seat_id, cinema_seat_id, status, price, show_id FROM show_seat WHERE show_id = $1`

	rows, err := psql.DB.Query(stmt, showID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all show seats for admin page: %w", err)
	}

	defer rows.Close()

	var allShowSeats []ShowSeatForAdmin

	for rows.Next() {
		var showSeat ShowSeatForAdmin

		err := rows.Scan(&showSeat.ShowSeatID, &showSeat.CinemaSeatID, &showSeat.SeatStatus, &showSeat.SeatPrice, &showSeat.ShowID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrShowSeatNotFound
			}
			return nil, fmt.Errorf("failed to scan show seat: %w", err)
		}

		allShowSeats = append(allShowSeats, showSeat)
	}

	return allShowSeats, nil
}

func (psql *Postgres) UpdateShowSeatByID(seatPrice int, showSeatID int) error {
	stmt := `UPDATE show_seat SET price = $1 WHERE show_seat_id = $2`

	result, err := psql.DB.Exec(stmt, seatPrice, showSeatID)
	if err != nil {
		return fmt.Errorf("error occurred while updating seat price: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error occurred while checking affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return ErrShowSeatNotFound
	}

	return nil
}
