package gormdb

import (
	"fmt"
	"log"
	"movie/models"
	"time"

	"gorm.io/gorm"
)

// PostgresDBRepo is the struct used to wrap our database connection pool, so that we
// can easily swap out a real database for a test database, or move to another database
// entirely, as long as the thing being swapped implements all of the functions in the type
// repository.DatabaseRepo.
type PostgresDBRepo struct {
	DB *gorm.DB
}

// Connection returns underlying connection pool.
func (m *PostgresDBRepo) Connection() *gorm.DB {
	return m.DB
}

func (m *PostgresDBRepo) AllMovies() ([]*models.Movie, error) {
	var movies []*models.Movie
	log.Println("movies")
	result := m.DB.Table("movies").Find(&movies)

	if result.Error != nil {
		return nil, result.Error
	}

	return movies, nil
}

func (m *PostgresDBRepo) MoviesByGenreID(id int) ([]*models.Movie, error) {
	var movies []*models.Movie

	where := ""
	if id > 0 {
		where = fmt.Sprintf("where id in (select movie_id from movies_genres where genre_id = %d)", id)
	}

	query := fmt.Sprintf(`
		select
			id, title, release_date, runtime,
			mpaa_rating, description, coalesce(image, ''),
			created_at, updated_at
		from
			movies %s
		order by
			title
	`, where)

	result := m.DB.Table("movies").Raw(query).Scan(&movies)

	if result.Error != nil {
		return nil, result.Error
	}
	return movies, nil
}

func (m *PostgresDBRepo) GetMovieGenresByMovieID(id int) (map[uint]struct{}, error) {
	genresIDs := map[uint]struct{}{}
	query := "select genre_id from movies_genres where movie_id=?"
	rows, err := m.DB.Table("movies").Raw(query, id).Rows()
	if err != nil {
		return nil, err
	}
	var genreID uint
	for rows.Next() {
		rows.Scan(&genreID)
		genresIDs[genreID] = struct{}{}
	}

	defer rows.Close()

	return genresIDs, nil
}

func (m *PostgresDBRepo) AllGenres() ([]*models.Genre, error) {
	var genres []*models.Genre
	result := m.DB.Table("genres").Find(&genres)

	if result.Error != nil {
		return nil, result.Error
	}

	return genres, nil
}

func (m *PostgresDBRepo) GetAMovie(id int) (*models.Movie, error) {
	var movie models.Movie
	result := m.DB.Table("movies").Find(&movie, id)

	if result.Error != nil {
		return nil, result.Error
	}

	return &movie, nil
}

func (m *PostgresDBRepo) AddMovie(movie *models.Movie, genres []*models.Genre) (int, error) {

	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()
	result := m.DB.Table("movies").Create(movie)
	if result.Error != nil {
		return 0, result.Error
	}

	//add genres by movie id
	for _, genre := range genres {
		if genre.Checked {
			movieGenre := models.MoviewGenre{
				MovieID: movie.ID,
				GenreID: genre.ID,
			}
			result := m.DB.Table("movies_genres").Create(&movieGenre)
			if result.Error != nil {
				return 0, result.Error
			}
		}
	}

	return int(result.RowsAffected), nil
}

func (m *PostgresDBRepo) UpdateMovie(movie *models.Movie) (int, error) {
	movie.UpdatedAt = time.Now()
	result := m.DB.Table("movies").Omit("created_at").Updates(movie)

	if result.Error != nil {
		return 0, result.Error
	}

	return int(result.RowsAffected), nil
}

func (m *PostgresDBRepo) UpdateMovieGenres(movieID uint, genres []*models.Genre) error {
	//delete genres by movie id
	m.DB.Table("movies_genres").Where("movie_id=?", movieID).Delete(models.MoviewGenre{})

	for _, genre := range genres {
		if genre.Checked {
			movieGenre := models.MoviewGenre{
				MovieID: movieID,
				GenreID: genre.ID,
			}
			result := m.DB.Table("movies_genres").Create(&movieGenre)
			if result.Error != nil {
				return result.Error
			}
		}
	}

	return nil
}

// DeleteMovie deletes one movie, by id.
func (m *PostgresDBRepo) DeleteMovie(id int) error {

	result := m.DB.Table("movies").Delete(&models.Movie{}, id)
	if result.Error != nil {
		return result.Error
	}
	m.DB.Table("movies_genres").Where("movie_id=?", id).Delete(models.MoviewGenre{})

	return nil
}
