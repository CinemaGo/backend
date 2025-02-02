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
