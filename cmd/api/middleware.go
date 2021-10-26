package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"

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

	// Declare a mutex and a map to hold the clients IP addresses and rate limiters.
	var (
		mu      sync.Mutex
		clients = make(map[string]*rate.Limiter)
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the client's IP address from the request
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// Lock the mutex to prevent this code from being executed
		mu.Lock()

		// Check the see if the IP address already exists in the map. If it doesn't, then initialize a new rate limiter and add the IP address and limiter to the map.
		if _, found := clients[ip]; !found {
			clients[ip] = rate.NewLimiter(2, 4)
		}

		// Call the Allow() method on the rate limiter for the current IP address. If the request isn't allowed, unlock the mutex and send a 429 Too many Requests response, just like before.
		if !clients[ip].Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		// Very importantly, unlock the mutex before calling the next handler in the chain.
		// Notice that we DON't use defer to unlock the mutex, as that would mean that the mutex isn't unlocked until all the handlers downstream of this middleware have also returned
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
