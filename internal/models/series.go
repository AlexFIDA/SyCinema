package models

type Series struct {
	ID          int
	Title       string
	Description string
	PosterURL   string
	ReleaseYear int
}

type Episode struct {
	ID         int
	SeriesID   int
	Season     int
	EpisodeNum int
	Title      string
	VideoPath  string
	IsFavorite bool
}