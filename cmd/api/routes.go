package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// initialize a new httprouter router instance
	router := httprouter.New()

	// overwritting the httprouter
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// register the relevant methods
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.getMoviesHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.getMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))

	// Users endpoint
	router.HandlerFunc(http.MethodPost, "/v1/users/register", app.registerUserHandler)

	// Activation endpoint
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	// authentication endpoint ==> /v1/login
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	// Wrap the router with the recovery middleware.
	// Use the authenticate() middleware on all requests.
	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
