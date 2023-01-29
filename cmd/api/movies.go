package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/NPeykov/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "creating a new movie")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDParam(r)
    if err != nil {
        http.NotFound(w, r) 
        return
    }

    movie := data.Movie {
        ID: id,
        CreatedAt: time.Now(),
        Title: "Casablanca",
        Runtime: 102,
        Genres: []string{"drama", "romance"},
        Version: 1,
    }

    err = app.writeJSON(w, http.StatusOK, movie, nil)

    if err != nil {
        app.logger.Println(err)
        http.Error(w, "The server can't process the request", http.StatusInternalServerError)
    }
}
