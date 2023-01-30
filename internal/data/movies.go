package data

import (
	"time"

	"github.com/NPeykov/greenlight/internal/validator"
)

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
