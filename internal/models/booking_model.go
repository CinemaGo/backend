package models

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

type DBContractBooking interface {
	RetrieveShowMovieInfo(showID int) (ShowMovieInfo, error)
	RetrieveShowInfo(movieID int) ([]ShowInfo, error)
	RetrieveShowStartTimes(showDate string) ([]ShowStartTime, error)
	RetrieveShowSeats(showID int) ([]ShowSeat, error)
	RetrieveShowSeatsMovieInfo(showID int) (ShowSeatsMovieInfo, error)

	RetrieveShowSeatStatus(showSeatID, showID int) (string, error)
	UpdateShowSeatStatus(status string, showSeatID, showID int) error
	InsertNewBooking(numberOfSeats int, paymentStatus string, userID, showID int) (int, error)
	InsertPaymentDetails(amount, remoteTransactionID int, paymentMethod string, bookingID int) error
}

type Bookings struct {
	db DBContractBooking
}

func NewBooking(db DBContractBooking) *Bookings {
	return &Bookings{db: db}
}

// RetrieveShowMovieInfo gets detailed information about the movie for a given show ID.
// This includes things like movie title, genre, age limit, and language. Essentially,
// it fetches the info you need for the movie that's playing in a particular show.
//
// Params:
//   - showID (int): The ID of the show for which we need to retrieve movie details.
//
// Returns:
//   - ShowMovieInfo: A struct containing movie details (ID, title, genre, age limit, language).
//   - error: Returns an error if something goes wrong while querying the database.
func (psql *Postgres) RetrieveShowMovieInfo(showID int) (ShowMovieInfo, error) {
	stmt := `SELECT m.id AS movie_id, m.title AS movie_title, m.genre AS movie_genre, m.age_limit AS movie_age_limit, m.language AS movie_language FROM movies m JOIN show s ON m.id = s.movie_id WHERE s.show_id = $1`

	var showMovieInfo ShowMovieInfo

	// Execute the query and map the result into our struct.
	err := psql.DB.QueryRow(stmt, showID).Scan(&showMovieInfo.MovieID, &showMovieInfo.MovieTitle, &showMovieInfo.MovieGenre, &showMovieInfo.MovieAgeLimit, &showMovieInfo.MovieLanguage)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// If no movie is found for the show ID, return a custom error.
			return ShowMovieInfo{}, ErrShowNotFound
		}

		// Any other error will be returned here with a helpful message.
		return ShowMovieInfo{}, fmt.Errorf("failed to retrieve movie show details for booking: %w", err)
	}

	// Return the movie details we fetched.
	return showMovieInfo, nil
}

// RetrieveShowInfo retrieves a list of cinema halls along with the corresponding show dates
// for a given movie. It queries the database to get unique combinations of hall names,
// hall types, and show dates for a specific movie ID, and orders the results by show date.
//
// Params:
//   - movieID (int): The ID of the movie for which we need to fetch show information.
//
// Returns:
//   - []ShowInfo: A slice of ShowInfo structs, each containing hall name, hall type,
//     and show date for the movie.
//   - error: If something goes wrong during the database query or data processing, an error
//     will be returned.
func (psql *Postgres) RetrieveShowInfo(movieID int) ([]ShowInfo, error) {
	stmt := `SELECT DISTINCT ch.hall_name, ch.hall_type, s.show_date FROM cinema_hall ch JOIN show s ON ch.cinema_hall_id = s.hall_id WHERE s.movie_id = $1 ORDER BY s.show_date`

	// Execute the query and get the rows of data.
	rows, err := psql.DB.Query(stmt, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve show infos form th database: %w", err)
	}

	defer rows.Close()

	var showInfos []ShowInfo

	// Loop through the rows and scan the data into our struct.
	for rows.Next() {
		var showInfo ShowInfo

		err := rows.Scan(&showInfo.HallName, &showInfo.HallType, &showInfo.ShowDate)
		if err != nil {
			// Handle edge cases where no rows are returned (not necessarily an error).
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrShowNotFound
			}
			// Return any other scan errors.
			return nil, fmt.Errorf("failed to scan show info: %w", err)
		}
		// Append the scanned show info to our slice.
		showInfos = append(showInfos, showInfo)
	}
	// Return the filled slice of show information.
	return showInfos, nil
}

// RetrieveShowStartTimes retrieves a list of show IDs and their corresponding start times
// for a given show date. It queries the database to fetch all shows for the specified
// date and returns the start times for each show.
//
// Params:
//   - showDate (string): The date for which we need to fetch the start times of the shows.
//
// Returns:
//   - []ShowStartTime: A slice of ShowStartTime structs, each containing a show ID and
//     its corresponding start time for the given date.
//   - error: If something goes wrong during the database query or data scanning,
//     an error will be returned.
func (psql *Postgres) RetrieveShowStartTimes(showDate string) ([]ShowStartTime, error) {
	stmt := `SELECT show_id, start_time FROM show WHERE show_date = $1`

	// Execute the query to fetch the show ID and start time for the given show date.
	rows, err := psql.DB.Query(stmt, showDate)
	if err != nil {
		// Return an error if the query fails.
		return nil, fmt.Errorf("failed to retrieve show start times from the database: %w", err)
	}

	// Ensure that rows are closed after processing to avoid resource leaks.
	defer rows.Close()

	var showStartTimes []ShowStartTime

	// Iterate over the rows and extract the start times.
	for rows.Next() {
		var startTime ShowStartTime

		// Scan the result into the ShowStartTime struct.
		err := rows.Scan(&startTime.ShowID, &startTime.StartTime)
		if err != nil {
			// If no rows are found, return a custom error.
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrStartTimeNotFound
			}
			// Return any other scanning errors.
			return nil, fmt.Errorf("failed to scan show start time: %w", err)
		}

		// Append the successfully scanned start time to the result slice.
		showStartTimes = append(showStartTimes, startTime)
	}

	// Return the list of show start times.
	return showStartTimes, nil
}

// RetrieveShowSeats retrieves a list of all seats for a given show, including seat details
// such as the row, seat number, type, show seat ID, status, and price. It queries the database
// to get this information by joining the `cinema_seat`, `show_seat`, and `show` tables.
//
// Params:
//   - showID (int): The ID of the show for which the seats need to be retrieved.
//
// Returns:
//   - []ShowSeat: A slice of ShowSeat structs containing information about each seat
//     for the specified show, including row, seat number, type, status, and price.
//   - error: An error if the query fails, or if there is any issue scanning the results.
func (psql *Postgres) RetrieveShowSeats(showID int) ([]ShowSeat, error) {
	stmt := `SELECT cs.seat_row, cs.seat_number, cs.seat_type, ss.show_seat_id, ss.status, ss.price FROM cinema_seat cs JOIN show_seat ss ON cs.cinema_seat_id = ss.cinema_seat_id JOIN show s ON ss.show_id = s.show_id WHERE s.show_id = $1;`

	// Execute the query using the provided showID.
	rows, err := psql.DB.Query(stmt, showID)
	if err != nil {
		// Return an error if the query fails.
		return nil, fmt.Errorf("failed to retrieve show seats from the database: %w", err)
	}

	// Ensure that rows are closed after processing to avoid resource leaks.
	defer rows.Close()

	var showSeats []ShowSeat

	// Iterate through the result rows and scan each seat's information into the ShowSeat struct.
	for rows.Next() {
		var showSeat ShowSeat

		// Scan the current row of data into the showSeat struct.
		err := rows.Scan(&showSeat.SeatRow, &showSeat.SeatNumber, &showSeat.SeatType,
			&showSeat.ShowSeatID, &showSeat.SeatStatus, &showSeat.SeatPrice)
		if err != nil {
			// Handle the case where no rows are found.
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrShowSeatNotFound
			}
			// Return an error if scanning fails.
			return nil, fmt.Errorf("failed to scan show seats: %w", err)
		}

		// Append the successfully scanned seat to the result slice.
		showSeats = append(showSeats, showSeat)
	}

	// Return the list of show seats.
	return showSeats, nil
}

// RetrieveShowSeatsMovieInfo retrieves movie details (such as title, show ID, show date, and start time)
// for a specific show using the showID. This is used to gather necessary information for the seats page.
//
// Params:
//   - showID (int): The ID of the show for which the movie information needs to be retrieved.
//
// Returns:
//   - ShowSeatsMovieInfo: A struct containing movie title, show ID, show date, and start time.
//   - error: An error if the query fails or if there is an issue retrieving or scanning the results.
func (psql *Postgres) RetrieveShowSeatsMovieInfo(showID int) (ShowSeatsMovieInfo, error) {
	stmt := `SELECT m.title AS movie_title, s.show_id, s.show_date, s.start_time FROM show s JOIN movies m ON s.movie_id = m.id WHERE s.show_id = $1`

	// Define a variable to hold the result.
	var showSeatsMovieInfo ShowSeatsMovieInfo

	// Execute the query and scan the results into the showSeatsMovieInfo struct.
	err := psql.DB.QueryRow(stmt, showID).Scan(&showSeatsMovieInfo.MovieTitle, &showSeatsMovieInfo.ShowID,
		&showSeatsMovieInfo.ShowDate, &showSeatsMovieInfo.ShowStartTime)
	if err != nil {
		// Handle the error if no results are found for the given showID.
		if errors.Is(err, ErrShowNotFound) {
			return ShowSeatsMovieInfo{}, ErrShowNotFound
		}
		// Return an error with more context if something goes wrong.
		return ShowSeatsMovieInfo{}, fmt.Errorf("failed to retrieve movie title, show date, show start time for seats page: %w", err)
	}

	// Return the populated struct containing movie and show details.
	return showSeatsMovieInfo, nil
}

// RetrieveShowSeatStatus retrieves the status of a specific seat for a given show.
// The status is returned as a string (e.g., "Available", "Selected", etc.).
//
// Params:
//   - showSeatID (int): The ID of the seat to retrieve the status for.
//   - showID (int): The ID of the show to which the seat belongs.
//
// Returns:
//   - string: The status of the seat, such as "Available", "Selected", etc.
//   - error: An error if the query fails, or if the seat or show is not found.
func (psql *Postgres) RetrieveShowSeatStatus(showSeatID, showID int) (string, error) {
	// SQL query to retrieve the status of a specific seat for a given show.
	stmt := `SELECT status FROM show_seat WHERE show_seat_id = $1 AND show_id = $2`

	// Variable to hold the seat status.
	var showSeatStatus string

	// Execute the query and scan the result into the showSeatStatus variable.
	err := psql.DB.QueryRow(stmt, showSeatID, showID).Scan(&showSeatStatus)
	if err != nil {
		// Handle case where no result is found for the given showSeatID and showID.
		if errors.Is(err, sql.ErrNoRows) {
			// Return a specific error indicating no data found.
			return "", ErrShowNotFound
		}
		// Return a formatted error with more context if any other error occurs.
		return "", fmt.Errorf("failed to retrieve show seat status: %w", err)
	}

	// Return the status of the seat if no errors occurred.
	return showSeatStatus, nil
}

// UpdateShowSeatStatus updates the status of a specific seat for a given show.
//
// This function allows updating the status of a seat, such as marking it as "Selected", "Booked", etc.
//
// Params:
//   - status (string): The new status to assign to the seat (e.g., "Selected", "Booked").
//   - showSeatID (int): The ID of the seat to update the status for.
//   - showID (int): The ID of the show to which the seat belongs.
//
// Returns:
//   - error: An error if the query fails, or if the seat/show combination is not found.
func (psql *Postgres) UpdateShowSeatStatus(status string, showSeatID, showID int) error {
	// SQL query to update the status of a specific seat for a given show.
	stmt := `UPDATE show_seat SET status = $1 WHERE show_seat_id = $2 AND show_id = $3`

	// Execute the update query with the provided status, showSeatID, and showID.
	_, err := psql.DB.Exec(stmt, status, showSeatID, showID)
	if err != nil {
		// Handle case where no rows are affected by the update (i.e., seat/show combination not found).
		if errors.Is(err, sql.ErrNoRows) {
			// Return a specific error indicating no data was found to update.
			return ErrShowNotFound
		}
		// Return a formatted error with more context if any other error occurs.
		return fmt.Errorf("failed to update show seat status: %w", err)
	}

	// Return nil if the status update was successful.
	return nil
}

// InsertNewBooking creates a new booking entry in the database.
//
// This function inserts a new booking record with the number of seats, payment status, user ID, and show ID.
//
// Params:
//   - numberOfSeats (int): The number of seats to be booked.
//   - paymentStatus (string): The current payment status (e.g., "Pending").
//   - userID (int): The ID of the user making the booking.
//   - showID (int): The ID of the show for which the seats are being booked.
//
// Returns:
//   - bookingID (int): The ID of the newly created booking.
//   - error: An error if the insertion fails.
func (psql *Postgres) InsertNewBooking(numberOfSeats int, paymentStatus string, userID, showID int) (int, error) {
	// SQL query to insert a new booking into the 'booking' table and return the booking ID.
	stmt := `INSERT INTO booking (number_of_seats, status, user_id, show_id) VALUES ($1, $2, $3, $4) RETURNING booking_id`

	var bookingID int
	// Execute the query and retrieve the generated booking ID.
	err := psql.DB.QueryRow(stmt, numberOfSeats, paymentStatus, userID, showID).Scan(&bookingID)
	if err != nil {
		// Return an error if the insertion fails.
		return 0, fmt.Errorf("failed to insert new booking into the database: %w", err)
	}

	// Return the generated booking ID if the insertion is successful.
	return bookingID, nil
}

// InsertPaymentDetails inserts payment details for a specific booking into the database.
//
// This function records the payment information for a booking, including the amount,
// remote transaction ID, payment method, and the associated booking ID.
//
// Params:
//   - amount (int): The payment amount for the booking.
//   - remoteTransactionID (int): The unique transaction ID from a remote payment gateway.
//   - paymentMethod (string): The method used for the payment (e.g., "Credit Card", "PayPal").
//   - bookingID (int): The ID of the booking for which the payment is being recorded.
//
// Returns:
//   - error: An error if the insertion fails, otherwise nil.
func (psql *Postgres) InsertPaymentDetails(amount, remoteTransactionID int, paymentMethod string, bookingID int) error {
	// SQL query to insert payment details into the 'payment' table.
	stmt := `INSERT INTO payment (amount, remote_transaction_id, payment_method, booking_id) VALUES ($1, $2, $3, $4)`

	// Execute the query to insert the payment details.
	_, err := psql.DB.Exec(stmt, amount, remoteTransactionID, paymentMethod, bookingID)
	if err != nil {
		// Return an error if the insertion fails.
		return fmt.Errorf("failed to insert payment details into the database: %w", err)
	}

	// Return nil if the payment details are successfully inserted.
	return nil
}
