package data

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/wagbubu/greenlight/internal/validator"
)

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) GetAll() ([]*Movie, error) {
	stmt := `SELECT id, created_at, title, year, runtime, genres, version FROM movies`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies := []*Movie{}

	for rows.Next() {
		movie := &Movie{}

		err := rows.Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version)
		if err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (m MovieModel) Insert(movie *Movie) error {
	stmt := `INSERT INTO movies (title, year, runtime, genres) VALUES ($1, $2, $3, $4) RETURNING id, created_at, version`

	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.DB.QueryRow(stmt, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	stmt := `SELECT id, created_at, title, year, runtime, genres, version FROM movies WHERE id = $1`
	var movie Movie
	err := m.DB.QueryRow(stmt, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		return nil, err
	}
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	stmt := `UPDATE movies SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1 WHERE id = $5 RETURNING version`
	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	return m.DB.QueryRow(stmt, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int64) error {
	stmt := `DELETE FROM movies WHERE id = $1`

	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

type Movie struct {
	ID        int64      `json:"id,omitempty"`         // Unique integer ID for the movie
	CreatedAt *time.Time `json:"created_at,omitempty"` // Timestamp for when the movie is added to our database
	Title     string     `json:"title,omitempty"`      // Movie title
	Year      int32      `json:"year,omitempty"`       // Movie release year
	Runtime   Runtime    `json:"runtime,omitempty"`    // Movie runtime (in minutes)
	Genres    []string   `json:"genres,omitempty"`     // Slice of genres for the movie (romance, comedy, etc.)
	Version   int32      `json:"version,omitempty"`    // The version number starts at 1 and will be incremented each
	UpdatedAt *time.Time `json:"updated_at,omitempty"` // time the movie information is updated
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}
