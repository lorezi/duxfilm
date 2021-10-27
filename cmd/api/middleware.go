package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

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

	// Define a client struct to hold the rate limiter and last seen time for each client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Declare a mutex and a map to hold the clients IP addresses and rate limiters.
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			// Lock the mutex to prevent any rate limiter checks from happening while the cleanup is taking place
			mu.Lock()

			// Loop through all clients. If they haven't been seen within the last three minutes, delete the corresponding entry from the map.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

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
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
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
