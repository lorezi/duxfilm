package main

import (
	"net/http"
)

// Declare a handler which writes a plain-text response with information about the application status, operating environment and version

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	// map to hold the server status info
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	// convert data map type to JSON
	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}
