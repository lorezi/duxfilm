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
