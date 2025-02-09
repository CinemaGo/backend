package handlers

import "time"

type newUserForm struct {
	Name            string `json:"name" binding:"required"`
	Surname         string `json:"surname" binding:"required"`
	Email           string `json:"email" binding:"required"`
	PhoneNumber     string `json:"phone_number" binding:"required"`
	Password        string `json:"password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}
type userLoginForm struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type userInfoUpdateFrom struct {
	Name        string `json:"name" binding:"required"`
	Surname     string `json:"surname" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
}

type BookingForm struct {
	ShowID      int   `json:"show_id" binding:"required"`
	ShowSeatsID []int `json:"show_seats_id" binding:"required"`
}

type NewCarouselImageForm struct {
	ImageURL      string `json:"carousel_image_image_url" binding:"required"`
	Title         string `json:"carousel_image_title" binding:"required"`
	Description   string `json:"carousel_image_description" binding:"required"`
	OrderPriority int    `json:"carousel_image_order_priority" binding:"required"`
}

type EditCarouselImageForm struct {
	CarouselImageID int    `json:"carousel_image_id" binding:"required"`
	ImageURL        string `json:"carousel_image_image_url" binding:"required"`
	Title           string `json:"carousel_image_title" binding:"required"`
	Description     string `json:"carousel_image_description" binding:"required"`
	OrderPriority   int    `json:"carousel_image_order_priority" binding:"required"`
}

type DeleteCarouselImageForm struct {
	CarouselImageID int `json:"carousel_image_id" binding:"required"`
}

type NewMovieForm struct {
	Title          string  `json:"title" binding:"required"`
	Description    string  `json:"description"  binding:"required"`
	Genre          string  `json:"genre"  binding:"required"`
	Language       string  `json:"language"  binding:"required"`
	TrailerURL     string  `json:"trailer_url"  binding:"required"`
	PosterURL      string  `json:"poster_url"  binding:"required"`
	Rating         float32 `json:"rating"  binding:"required"`
	RatingProvider string  `json:"rating_provider"  binding:"required"`
	Duration       int     `json:"duration"  binding:"required"`
	ReleaseDate    string  `json:"release_date"  binding:"required"`
	AgeLimit       string  `json:"age_limit"  binding:"required"`
}

type EditMovieForm struct {
	MovieID        int     `json:"movie_id" binding:"required"`
	Title          string  `json:"title" binding:"required"`
	Description    string  `json:"description"  binding:"required"`
	Genre          string  `json:"genre"  binding:"required"`
	Language       string  `json:"language"  binding:"required"`
	TrailerURL     string  `json:"trailer_url"  binding:"required"`
	PosterURL      string  `json:"poster_url"  binding:"required"`
	Rating         float32 `json:"rating"  binding:"required"`
	RatingProvider string  `json:"rating_provider"  binding:"required"`
	Duration       int     `json:"duration"  binding:"required"`
	ReleaseDate    string  `json:"release_date"  binding:"required"`
	AgeLimit       string  `json:"age_limit"  binding:"required"`
}

type DeleteMovieForm struct {
	MovieID int `json:"movie_id" binding:"required"`
}

type NewActorCrewForm struct {
	FullName        string    `json:"full_name" binding:"required"`
	ImageURL        string    `json:"image_url" binding:"required"`
	Occupation      string    `json:"occupation" binding:"required"`
	RoleDescription string    `json:"role_description" binding:"required"`
	BornDate        time.Time `json:"born_date" binding:"required"`
	Birthplace      string    `json:"birthplace" binding:"required"`
	About           string    `json:"about" binding:"required"`
	IsActor         bool      `json:"is_actor" binding:"required"`
	MovieID         int       `json:"movie_id" binding:"required"`
}

type EditActorCrewForm struct {
	ActorCrewID     int       `json:"actor_crew_id" binding:"required"`
	FullName        string    `json:"full_name" binding:"required"`
	ImageURL        string    `json:"image_url" binding:"required"`
	Occupation      string    `json:"occupation" binding:"required"`
	RoleDescription string    `json:"role_description" binding:"required"`
	BornDate        time.Time `json:"born_date" binding:"required"`
	Birthplace      string    `json:"birthplace" binding:"required"`
	About           string    `json:"about" binding:"required"`
	IsActor         bool      `json:"is_actor" binding:"required"`
}

type DeleteActorCrewForm struct {
	ActorCrewID int `json:"actor_crew_id" binding:"required"`
}

type NewCinemaHallForm struct {
	HallName string `json:"hall_name" binding:"required"`
	HallType string `json:"hall_type" binding:"required"`
	Capacity int    `json:"capacity" binding:"required"`
}

type EditCinemaHallForm struct {
	CinemaHallID int    `json:"cinema_hall_id" binding:"required"`
	HallName     string `json:"hall_name" binding:"required"`
	HallType     string `json:"hall_type" binding:"required"`
	Capacity     int    `json:"capacity" binding:"required"`
}

type DeleteCinemaHallForm struct {
	CinemaHallID int `json:"cinema_hall_id" binding:"required"`
}

type NewCinemaSeatForm struct {
	SeatRow    string `json:"seat_row" binding:"required"`
	SeatNumber int    `json:"seat_number" binding:"required"`
	SeatType   string `json:"seat_type" binding:"required"`
	HallID     int    `json:"hall_id" binding:"required"`
}

type DeleteCinemaSeatForm struct {
	CinemaSeatID int `json:"cinema_seat_id" binding:"required"`
}

type NewShowForm struct {
	ShowDate  time.Time `json:"show_date" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	HallID    int       `json:"hall_id" binding:"required"`
	MovieID   int       `json:"movie_id" binding:"required"`
}
