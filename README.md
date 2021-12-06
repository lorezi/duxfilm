# Duxfilm 🎦

a JSON API for retrieving and managing information about movies.

# Features

1. RESTful API
2. Error management - Triaged the Decode error
   - Triaging the Decode error gracefully handles errors and makes the error messages simpler, clearer, consistent in formatting and avoids unnecessary exposure of information.
3. Custom request body validations
4. SQL database schema migrations
5. Simple CRUD operations
6. Advanced CRUD operations
   - Supports partial updates to a resource
   - Optimistic concurrency control to avoid race conditions when two clients tries to update the same resource at the same time
   - Context timeouts
7. Filtering, Pagination and Sorting
   - Reductive filtering - which allows clients to search based on a case insensitive exact match for movie.
   - Postgres full-text search functionality
   - Sorting lists
8. Rich structured logging and error handling
9. API rate limiting
   - Learn the principle behind token-bucket rate-limiter algorithms and how to apply them in context of an API or web application.
   - Use middleware to rate-limit requests to your API endpoints, first by making a single rate global limiter, then extending it to support per-client limiting based on IP address.
   - Configuration settings for rate limiter
10. Graceful shutdown
11. User model setup and registration
12. Email setup
    - Use `goroutine` to send email at the background
13. Recovering panics while using goroutine
14. Graceful shutdown of background tasks using syncWaitGroup https://play.golang.org/p/1j9-p8JOKWa
15. User Activation

- Activation functionality to confirm that a user used their own, real, email address by including 'account activation' instructions in their welcome email.
- Implement a secure 'account activation' workflow which verifies a new user's email address.
- Generate cryptographically secure random tokens
- Generate fast hashes of data
- Implement patterns for working with cross-table relationships in your database, including setting up foreign keys and retrieving related data via SQL JOIN queries
- Implements user activation after user account creation

16. Authentication feature (Know the user)

- Get authentication token
- Authentication requests

  - If token is valid, look up the user details and add their details to the request context.
  - If no Authorization header was provided at all, then we will add the details for an anonymous user to the request context instead.

17. Permission-based Authorization

- Add checks so that only activated users are able to access the various endpoints
- Implement a permission-based authorization pattern, which provides fine-grained control over exactly which users can access which endpoints
- Create a subset of permissions that only users who have a specific permission can perform specific operations.

18. CORS

- Supporting multiple dynamic origins
- Middleware to intercepts and responds to any preflight requests.

19. Metrics - Performance insight into your application - Instrumentation of the Golang app

- Gives insight to how much memory is the application using and how this is changing over time
- How many goroutines are currently in use and changes over time
- How many database connections are in use and how many are idle
- What's the ratio of successful HTTP responses to both client and server errors
- Use the expvar package to view application metrics in JSON format via a HTTP handler.
- Request-level metrics - Total number of requests received & responses sent and the total time taken to process all requests in microseconds

- UnMarshal json to golang struct: https://go.dev/play/p/by_OhVyTANB
- UnMarshal json to golang struct: https://go.dev/play/p/x3kEisafNPo

20. Use a makefile to automate common tasks in your project, such as creating and executing migrations.

- Carry out quality control checks of your code using the `go go vet` and staticcheck tools
- Generate version numbers automatically based on Git commits and integrate them into your application.
