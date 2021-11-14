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
