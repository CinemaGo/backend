package services

import (
	"cinemaGo/backend/internal/models"
	"errors"
	"fmt"
	"time"
)

type BookingServiceInterface interface {
	FetchShowMovieInfo(showID int) (models.ShowMovieInfo, error)
	FetchShowInfo(movieID int) ([]models.ShowInfo, error)
	FetchShowStartTimes(showDate string) ([]models.ShowStartTime, error)
	FetchShowSeats(showID int) ([]models.ShowSeat, error)
	FetchShowSeatsMovieInfo(showID int) (models.ShowSeatsMovieInfo, error)
	CreateNewBooking(showID, userID int, showSeatsID []int) error
}

type BookingService struct {
	db models.DBContractBooking
}

func NewBookingService(db models.DBContractBooking) *BookingService {
	return &BookingService{db: db}
}

// FetchShowMovieInfo fetches the movie details for a specific show.
//
// This function retrieves detailed information about a movie, such as its title, genre,
// age limit, and language, based on the provided show ID. It then returns the movie information.
//
// Params:
//   - showID (int): The ID of the show for which movie information is being fetched.
//
// Returns:
//   - models.ShowMovieInfo: A structure containing the movie's details (ID, title, genre, etc.)
//   - error: An error if the retrieval fails, otherwise nil.
func (bs *BookingService) FetchShowMovieInfo(showID int) (models.ShowMovieInfo, error) {
	// Fetch movie details from the database.
	showMovieInfo, err := bs.db.RetrieveShowMovieInfo(showID)
	if err != nil {
		// Handle case where the show is not found.
		if errors.Is(err, models.ErrShowNotFound) {
			return models.ShowMovieInfo{}, ErrShowNotFound
		}
		// Return a wrapped error if there is any other failure.
		return models.ShowMovieInfo{}, fmt.Errorf("error occurred while fetching show movie details in the service section: %w", err)
	}

	// Return the movie details if the fetch was successful.
	return showMovieInfo, nil
}

// FetchShowInfo retrieves the show information for a specific movie, including show details like
// the show date and associated start times for each show.
//
// This function fetches the show information (such as the hall name, show date, and hall type),
// formats the show date, and retrieves the corresponding show start times for the given movie ID.
//
// Params:
//   - movieID (int): The ID of the movie for which show information is being fetched.
//
// Returns:
//   - []models.ShowInfo: A list of show details for the specified movie.
//   - error: An error if the retrieval or processing fails, otherwise nil.
func (bs *BookingService) FetchShowInfo(movieID int) ([]models.ShowInfo, error) {
	// Fetch show information for the given movie from the database.
	showInfos, err := bs.db.RetrieveShowInfo(movieID)
	if err != nil {
		// Handle case where no shows were found for the movie.
		if errors.Is(err, models.ErrShowNotFound) {
			return nil, ErrShowNotFound
		}
		// Return a wrapped error if there is any failure in fetching show info.
		return nil, fmt.Errorf("error occurred while fetching show info in the service section: %w", err)
	}

	// Iterate over each show information record.
	for i := range showInfos {
		// Parse and format the show date.
		t, err := time.Parse(time.RFC3339, showInfos[i].ShowDate)
		if err != nil {
			return nil, fmt.Errorf("error occurred while formatting show date: %w", err)
		}
		showInfos[i].ShowDate = t.Format("2006-01-02")

		// Fetch corresponding show start times for the formatted show date.
		startTimes, err := bs.FetchShowStartTimes(showInfos[i].ShowDate)
		if err != nil {
			return nil, fmt.Errorf("error occurred while fetching show start times: %w", err)
		}

		// Format the start times for each show.
		for j := range startTimes {
			startTime, err := time.Parse(time.RFC3339, startTimes[j].StartTime)
			if err != nil {
				return nil, fmt.Errorf("error occurred while formatting start time: %w", err)
			}
			startTimes[j].StartTime = startTime.Format("15:04")
		}

		// Assign the formatted start times to the show information.
		showInfos[i].ShowStartTimes = startTimes
	}

	// Return the list of show information, including start times, after processing.
	return showInfos, nil
}

// FetchShowStartTimes retrieves the start times for shows on a specific date.
//
// This function fetches the list of start times for shows that are scheduled on a particular date.
// It returns the list of start times for shows on that date or an error if the retrieval fails.
//
// Params:
//   - showDate (string): The date of the shows for which start times are being fetched.
//
// Returns:
//   - []models.ShowStartTime: A list of start times for shows on the specified date.
//   - error: An error if the retrieval fails, otherwise nil.
func (bs *BookingService) FetchShowStartTimes(showDate string) ([]models.ShowStartTime, error) {
	// Retrieve start times for the given show date from the database.
	showStartTimes, err := bs.db.RetrieveShowStartTimes(showDate)
	if err != nil {
		// Handle case where no start times were found for the show date.
		if errors.Is(err, models.ErrStartTimeNotFound) {
			return nil, ErrStartTimeNotFound
		}
		// Return a wrapped error if there is any other failure in fetching start times.
		return nil, fmt.Errorf("error occurred while fetching show start times in the service section: %w", err)
	}

	// Return the list of show start times if the fetch was successful.
	return showStartTimes, nil
}

// FetchShowSeats retrieves the list of available seats for a specific show.
//
// This function fetches the details of all the seats available for a show, such as seat row,
// seat number, seat type, and the status (whether the seat is booked or available).
// It returns the list of seats for the specified show or an error if the retrieval fails.
//
// Params:
//   - showID (int): The ID of the show for which seat details are being fetched.
//
// Returns:
//   - []models.ShowSeat: A list of seat details for the specified show.
//   - error: An error if the retrieval fails, otherwise nil.
func (bs *BookingService) FetchShowSeats(showID int) ([]models.ShowSeat, error) {
	// Retrieve seat details for the given show ID from the database.
	showSeats, err := bs.db.RetrieveShowSeats(showID)
	if err != nil {
		// Handle case where no seats were found for the show ID.
		if errors.Is(err, models.ErrShowSeatNotFound) {
			return nil, ErrShowSeatNotFound
		}
		// Return a wrapped error if there is any other failure in fetching show seats.
		return nil, fmt.Errorf("error occurred while fetching show seats in the service section: %w", err)
	}

	// Return the list of show seats if the fetch was successful.
	return showSeats, nil
}

// FetchShowSeatsMovieInfo retrieves movie details, show date, and show start time for a specific show ID.
//
// This function fetches the movie title, show date, and show start time for a given show ID. It then
// formats the date and time for display purposes. The function returns the movie info along with
// formatted date and start time, or an error if the retrieval or formatting fails.
//
// Params:
//   - showID (int): The ID of the show for which movie details, date, and start time are being fetched.
//
// Returns:
//   - models.ShowSeatsMovieInfo: A structured object containing movie title, show date, and start time.
//   - error: An error if the retrieval or formatting fails, otherwise nil.
func (bs *BookingService) FetchShowSeatsMovieInfo(showID int) (models.ShowSeatsMovieInfo, error) {
	// Retrieve the movie details, show date, and show start time for the given show ID from the database.
	showSeatsMovieInfo, err := bs.db.RetrieveShowSeatsMovieInfo(showID)
	if err != nil {
		// If no show information is found, return the appropriate error.
		if errors.Is(err, models.ErrShowNotFound) {
			return models.ShowSeatsMovieInfo{}, ErrShowNotFound
		}
		// Wrap the error with additional context if the retrieval fails for any reason.
		return models.ShowSeatsMovieInfo{}, fmt.Errorf("error occurred while fetching movie title, show date, show start time for seats page in the service section: %w", err)
	}

	// Parse and format the show date.
	t, err := time.Parse(time.RFC3339, showSeatsMovieInfo.ShowDate)
	if err != nil {
		// Return an error if the date formatting fails.
		return models.ShowSeatsMovieInfo{}, fmt.Errorf("error occurred while formatting show date: %w", err)
	}
	// Update the ShowDate field with the formatted date.
	showSeatsMovieInfo.ShowDate = t.Format("2006-01-02")

	// Parse and format the show start time.
	startTime, err := time.Parse(time.RFC3339, showSeatsMovieInfo.ShowStartTime)
	if err != nil {
		// Return an error if the time formatting fails.
		return models.ShowSeatsMovieInfo{}, fmt.Errorf("error occurred while formatting start time: %w", err)
	}

	// Update the ShowStartTime field with the formatted start time.
	showSeatsMovieInfo.ShowStartTime = startTime.Format("15:04")

	// Return the movie info along with the formatted date and start time.
	return showSeatsMovieInfo, nil
}

// CreateNewBooking handles the creation of a new booking for a user by selecting seats and processing the booking.
//
// This function verifies the availability of the selected seats, updates their status, creates a new booking,
// updates the booking's seat status, and finally inserts the payment details into the database. It ensures that
// no more than five seats can be selected at once and handles any errors encountered during these operations.
//
// Params:
//   - showID (int): The ID of the show that the user is booking seats for.
//   - userID (int): The ID of the user who is making the booking.
//   - showSeatsID ([]int): A slice of seat IDs that the user is selecting for the booking.
//
// Returns:
//   - error: Returns nil if the booking was created successfully, or an error if any part of the process fails.
func (bs *BookingService) CreateNewBooking(showID, userID int, showSeatsID []int) error {
	var numberOfSeats int

	// Loop through the selected show seats to check if they are available.
	for _, showSeatID := range showSeatsID {
		// Retrieve the current status of the seat.
		showSeatStatus, err := bs.db.RetrieveShowSeatStatus(showSeatID, showID)
		if err != nil {
			// If an error occurs during status retrieval (show not found), return the error.
			if errors.Is(err, models.ErrShowNotFound) {
				return ErrShowNotFound
			}
			return fmt.Errorf("error occurred while fetching the status of the seat in the service section: %w", err)
		}

		// If the seat is not available, set the number of seats to 0 and return an error.
		if showSeatStatus != "Available" {
			numberOfSeats = 0
			return ErrShowSeatHasSelected
		}
		// Increment the number of seats selected.
		numberOfSeats++
	}

	// If more than 5 seats are selected, return an error.
	if numberOfSeats > 5 {
		return ErrTooManySeats
	}

	// Loop through the selected seats again to update their statuses and create the booking.
	for _, showSeatID := range showSeatsID {
		// Mark the seat as "Selected" in the database.
		err := bs.db.UpdateShowSeatStatus("Selected", showSeatID, showID)
		if err != nil {
			return fmt.Errorf("failed to update status of the show seat in the service section: %w", err)
		}

		// Create a new booking in the database with "Pending" payment status.
		bookingID, err := bs.db.InsertNewBooking(numberOfSeats, "Pending", userID, showID)
		if err != nil {
			return fmt.Errorf("error occurred while creating new booking in the service section: %w", err)
		}

		// Mark the seat as "Booked" after the booking is created.
		err = bs.db.UpdateShowSeatStatus("Booked", showSeatID, showID)
		if err != nil {
			return fmt.Errorf("failed to update status of the show seat in the service section: %w", err)
		}

		// Insert payment details with a "Pending" status (with a default amount of 100).
		err = bs.db.InsertPaymentDetails(100, 0, "", bookingID)
		if err != nil {
			return fmt.Errorf("error occurred while inserting payment details in the service section: %w", err)
		}
	}

	// Reset the number of seats after completing the booking process.
	numberOfSeats = 0

	// Return nil indicating the successful creation of the booking.
	return nil
}
