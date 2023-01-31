package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/NPeykov/greenlight/internal/validator"
	"github.com/lib/pq"
)

type MovieModel struct {
    DB *sql.DB
}

type Movie struct {
    ID int64            `json:"id"`
    CreatedAt time.Time `json:"-"`
    Title string        `json:"title"`
    Year int32          `json:"year,omitempty"`
    Runtime int32       `json:"runtime,omitempty"`
    Genres []string     `json:"genres,omitempty"`
    Version int32       `json:"version"`
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
    //check Title
    v.Check(movie.Title != "", "title", "must be provided")
    v.Check(len(movie.Title) <= 500, "title", "title has more than 500 bytes long")
    //check Year 
    v.Check(movie.Year != 0, "year", "must be provided")
    v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
    v.Check(movie.Year <= int32(time.Now().Year()), "title", "invalid year")
    //check Runtime 
    v.Check(movie.Runtime != 0, "runtime", "must be provided")
    v.Check(movie.Runtime > 0, "runtime", "must be a positive number")
    //check Genres 
    v.Check(movie.Genres != nil, "genres", "must be provided")
    v.Check(len(movie.Genres) > 0, "genres", "must have at least one gendre")
    v.Check(len(movie.Genres) <= 5, "genres", "must have less than five genres")
    v.Check(validator.Unique(movie.Genres), "genres", "genres cannot be duplicated")
}

func (m MovieModel) Insert(movie *Movie) error {
    query := `INSERT INTO movies (title, year, runtime, genres)
    VALUES ($1, $2, $3, $4)
    RETURNING id, created_at, version`
    args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
    return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {
    query := `SELECT id, created_at, title, year, runtime, genres, version
    FROM movies
    WHERE id = $1`

    var movie Movie

    err := m.DB.QueryRow(query, id).Scan(
        &movie.ID,
        &movie.CreatedAt,
        &movie.Title,
        &movie.Year,
        &movie.Runtime,
        pq.Array(&movie.Genres),
        &movie.Version,
    )

    if err != nil {
        switch {
        case errors.Is(err, sql.ErrNoRows):
            return nil, ErrRecordNotFound
        default:
            return nil, err
        }
    }

    return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
    query := `UPDATE movies
    SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
    WHERE id = $5
    RETURNING version`
    args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres), movie.ID}
    return m.DB.QueryRow(query, args...).Scan(&movie.Version)
}

func (m MovieModel) Delete(id int64) error {
    query := `DELETE FROM movies
    WHERE id = $1`

    result, err := m.DB.Exec(query, id)

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
