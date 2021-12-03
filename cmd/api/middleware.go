package main

import (
	"errors"
	"expvar"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/lorezi/duxfilm/internal/data"
	"github.com/lorezi/duxfilm/internal/validator"
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

		// Only carry out the check if rate limiting is enabled.
		if app.config.limiter.enabled {
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
				clients[ip] = &client{
					limiter: rate.NewLimiter(
						rate.Limit(app.config.limiter.rps), app.config.limiter.burst)}
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
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Authorization" header to the response. This indicates to any caches that the response may vary based on the value of the Authorization header in the request.
		w.Header().Add("Vary", "Authorization")

		// Retrieve the value of the Authorization header from the request. This will return the empty string "" if there is no such header found.
		authorizationHeader := r.Header.Get("Authorization")

		// If there is no Authorization header found, use the contextSetUser() helper  that we just made to add the AnonymousUser to the request context. Then we call the next handler in the chain and return without executing any of the code below.
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		// Otherwise, we expect the value of the Authorization header to be in the format "Bearer <token>". We try to split this into its constituent parts, and if the header isn't in the expected format we return a 401 Unauthorized response using the invalidAuthenticationTokenResponse() helper.
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Extract the actual authentication token from the header parts
		token := headerParts[1]

		// Validate the token to make sure it is in a sensible format.
		v := validator.New()

		// If the token isn't valid, use the invalidAuthenticationTokenResponse() helper to send a response, rather than the failedValidationResponse() helper that we'd normally use.
		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Retrieve the details of the user associated with the authentication token, again calling the invalidAuthenticationTokenResponse() helper if no matching record was found.
		// IMPORTANT: Notice that we are using ScopeAuthentication as the first parameter here.
		user, err := app.models.User.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidCredentialResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		// Call the contextSetUser() helper to add the user information to the request context.
		r = app.contextSetUser(r, user)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}

// RequireAuthenticatedUser() middleware to check that a user is not anonymous.
// This function would only check for known user
func (app *application) RequireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)

		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Checks that a user is both authenticated and activated
func (app *application) requireActivateUser(next http.HandlerFunc) http.HandlerFunc {
	// Assign the handlerFunc to a fn variable
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve user with contextGetUser()
		user := app.contextGetUser(r)

		// If the user is anonymous, then call the authenticationRequiredResponse() to inform the client that they should authenticate before trying again.
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		// Call the next handler in the chain.
		next.ServeHTTP(w, r)
	})

	// Wrap fn with the requireAuthenticatedUser() middleware before returning it.
	return app.RequireAuthenticatedUser(fn)
}

// Note that the first paramater for the middleware function is the permission code that we require the user to have.
func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn :=
		func(w http.ResponseWriter, r *http.Request) {
			// Retrieve the user from the request context.
			user := app.contextGetUser(r)

			// Get the slice of permissions for the user.
			permissions, err := app.models.Permission.GetAllForUser(user.ID)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			// Check if the slice includes the required permission. If it doesn't, then return a 403 Forbidden response.
			if !permissions.Include(code) {
				app.notPermittedResponse(w, r)
				return
			}

			// Otherwise they have the required permission so we call the next handler in the chain
			next.ServeHTTP(w, r)
		}
	return app.requireActivateUser(fn)
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Origin" header.
		w.Header().Add("Vary", "Origin")

		w.Header().Add("Vary", "Access-Control-Request-Method")

		// Get the value of the request's Origin header.
		origin := r.Header.Get("Origin")

		// Only run this if there's an Origin request header present AND at least one trusted origin is configured.
		if origin != "" && len(app.config.cors.trustedOrigins) != 0 {
			// Loop through the list of trusted origins, checking to see if the request origin exactly matches one of them.
			for i := range app.config.cors.trustedOrigins {
				if origin == app.config.cors.trustedOrigins[i] {
					// If there is a match, then set a "Access-Control-Allow-Origin"
					// response header with the request origin as the value.
					w.Header().Set("Access-Control-Allow-Origin", origin)

					// Check if the request has the HTTP method OPTIONS and contains the "Access-Control-Request-Method" header.
					// If it does, then we treat it as a preflight request.
					if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
						// Set the necessary preflight response headers, as discussed previously.
						w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						// Write the headers along with a 200 OK status and return from the middleware with no further action.
						w.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) metrics(next http.Handler) http.Handler {
	// Initialize the new expvar variables when the middleware chain in first built.
	totalRequestReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_ūs")

	// The following code will be run for every request...
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Record the time that we started to process the request.
		start := time.Now()

		// Use the Add() method to increment the number of requests received by 1.
		totalRequestReceived.Add(1)

		// Call the next handler in the chain
		next.ServeHTTP(w, r)

		// On the way back up the middleware chain, increment the number or responses sent by 1
		totalResponsesSent.Add(1)

		// Calculate the number of microseconds since we began to process the request,
		// then increment the total processing time by this amount.
		duration := time.Now().Sub(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(duration)
	})
}
