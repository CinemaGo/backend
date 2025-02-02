package models

type CarouselImage struct {
	ID       int
	ImageURL string
}

type AllMovies struct {
	ID             int
	Title          string
	Genre          string
	Language       string
	PosterUrl      string
	Rating         float32
	RatingProvider string
	AgeLimit       string
}

type Movie struct {
	ID             int
	Title          string
	Description    string
	Genre          string
	Language       string
	TrailerUrl     string
	PosterUrl      string
	Rating         float32
	RatingProvider string
	Duration       int
	ReleaseDate    string
	AgeLimit       string
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
