package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/wagbubu/greenlight/internal/data"
)

func (app *application) createNewMovieHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create movie handler")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDparam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	createdAt := time.Now()
	movie := data.Movie{
		ID:        id,
		CreatedAt: &createdAt,
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, movie, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "something went wrong!", http.StatusInternalServerError)
		return
	}
}
