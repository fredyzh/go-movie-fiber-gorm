package repositories

import (
	"movie/models"

	"gorm.io/gorm"
)

type DatabaseRepo interface {
	Connection() *gorm.DB
	AllMovies() ([]*models.Movie, error)
	AllGenres() ([]*models.Genre, error)
	GetAMovie(id int) (*models.Movie, error)
	AddMovie(*models.Movie, []*models.Genre) (int, error)
	MoviesByGenreID(id int) ([]*models.Movie, error)
	GetMovieGenresByMovieID(id int) ([]int, error)

	// OneMovieForEdit(id int) (*models.Movie, []*models.Genre, error)
	// OneMovie(id int) (*models.Movie, error)

	UpdateMovieGenres(movieID int, genres []*models.Genre) error
	UpdateMovie(movie *models.Movie) error
	DeleteMovie(id int) error
}
