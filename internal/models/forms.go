package models

import "time"

type CarouselImage struct {
	ID       int
	ImageURL string
}

type AllShowsMovie struct {
	ShowID              int
	MovieID             int
	MovieTitle          string
	MovieGenre          string
	MovieLanguage       string
	MoviePosterUrl      string
	MovieRating         float32
	MovieRatingProvider string
	MovieAgeLimit       string
}

type AShowMovie struct {
	ShowID              int
	MovieID             int
	MovieTitle          string
	MovieDescription    string
	MovieGenre          string
	MovieLanguage       string
	MovieTrailerUrl     string
	MoviePosterUrl      string
	MovieRating         float32
	MovieRatingProvider string
	MovieDuration       int
	MovieReleaseDate    string
	MovieAgeLimit       string
}

type ActorsCrewsOfMovie struct {
	ID              int
	FullName        string
	ImageURL        string
	RoleDescription string
	IsActor         bool
}

type ActorCrewInfo struct {
	ID              int
	FullName        string
	ImageURL        string
	Occupation      string
	RoleDescription string
	BornDate        string
	Birthplace      string
	About           string
}

type ActorCrewMovies struct {
	ID        int
	Title     string
	PosterUrl string
}

type UserInfo struct {
	Name        string
	Surname     string
	Email       string
	PhoneNumber string
}

type ShowMovieInfo struct {
	MovieID       int
	MovieTitle    string
	MovieGenre    string
	MovieAgeLimit string
	MovieLanguage string
}

type ShowInfo struct {
	HallName       string
	HallType       string
	ShowDate       string
	ShowStartTimes []ShowStartTime
}

type ShowStartTime struct {
	ShowID    int
	StartTime string
}

type ShowSeat struct {
	ShowSeatID int
	SeatRow    string
	SeatNumber int
	SeatType   string
	SeatStatus string
	SeatPrice  int
}

type ShowSeatsMovieInfo struct {
	MovieTitle    string
	ShowID        int
	ShowDate      string
	ShowStartTime string
}

type CarouselImageForAdmin struct {
	CarouselImageID int
	ImageURL        string
	Title           string
	Description     string
	OrderPriority   int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type AllMoviesForAdmin struct {
	ID    int
	Title string
}

type MovieForAdmin struct {
	MovieID        int
	Title          string
	Description    string
	Genere         string
	Language       string
	TrailerURL     string
	PosterURL      string
	Rating         float32
	RatingProvider string
	Duration       int
	RelaseDate     string
	AgeLimit       string
	CreatedAt      time.Time
	UpdatesAt      time.Time
}

type AllActorCrewForAdmin struct {
	ActorCrewID     int
	FullName        string
	ImageURL        string
	RoleDescription string
	IsActor         bool
}

type ActorCrewForAdmin struct {
	ActorCrewID     int
	FullName        string
	ImageURL        string
	Occupation      string
	RoleDescription string
	BornDate        string
	Birthplace      string
	About           string
	IsActor         bool
}

type CinemaHallForAdmin struct {
	CinemaHallID int
	HallName     string
	HallType     string
	Capacity     int
}

type CinemaSeatForAdmin struct {
	CinemaSeatID int
	SeatRow      *string
	SeatNumber   int
	SeatType     string
	HallID       int
}

type ShowForAdmin struct {
	ShowID    int
	ShowDate  string
	StartTime string
	HallID    int
	MovieID   int
}

type ShowSeatForAdmin struct {
	ShowSeatID   int
	CinemaSeatID int
	SeatStatus   string
	SeatPrice    float32
	ShowID       int
}
