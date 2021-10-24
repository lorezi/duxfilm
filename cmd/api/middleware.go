package main

import (
	"fmt"
	"net/http"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event of a panic or not)

		// middleware logic comes here...
		defer func() {
			if err := recover(); err != nil {
				// User the builtin recover function to check if there has been a panic or not.
				rw.Header().Set("Connection", "close")

				// The value returned by recover() has the type interface{}, so we use fmt.Errorf() to normalize it into error and call our serverErrorResponse() helper.

				// In turn, this will log the error using our custom Logger type at the ERROR level and send the client a 500 Internal Server Error Response
				app.serverErrorResponse(rw, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(rw, r)

	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	// Initialize a new rate limiter which allows an average of 2 requests per second, with a maximum of 4 requests in a single `burst`
	limiter := rate.NewLimiter(2, 4)

	// The function we are returning is a closure, which 'closes over' the limiter variable.
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Call limiter.Allow() to see if the request is permitted, and if it's not, then we call the rateLimitExceededResponse() helper to return a 429 Too Many Requests response(we will create this helper in a minute).
		if !limiter.Allow() {
			app.rateLimitExceededResponse(rw, r)
			return
		}

		next.ServeHTTP(rw, r)
	})
}
