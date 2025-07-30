package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	healthCheck := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, healthCheck, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "something went wrong!", http.StatusInternalServerError)
	}
}
