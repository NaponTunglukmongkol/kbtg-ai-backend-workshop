# System Prompt for GitHub Copilot

## Coding Standards
- Follow Go best practices for clean, readable, and maintainable code.
- Use descriptive variable and function names (e.g., `setupRoutes`, `initDatabase`).
- Ensure proper error handling and logging (e.g., `log.Fatal`, `c.Status`).
- Use consistent indentation (tabs for Go).
- Write comments for complex logic and SQL queries.

## High-Level Ideas
- Focus on modular and reusable code (e.g., separate database initialization and route setup).
- Prioritize performance and scalability (e.g., efficient SQL queries, proper indexing).
- Ensure security best practices:
  - Validate user inputs to prevent SQL injection.
  - Use prepared statements for database queries.
- Write unit tests for critical functionality (e.g., CRUD operations).
- Use Fiber framework features effectively for routing and middleware.

## Scope and Limitations
### Allowed
- Implementing CRUD operations for `users`, `transfers`, and `point_ledger`.
- Writing tests for database interactions and API endpoints.
- Refactoring code for better readability and performance.
- Updating Swagger documentation (`swagger.yml`) to reflect API changes.

### Not Allowed
- Writing code that violates ethical guidelines (e.g., unauthorized data access).
- Generating content that is harmful, hateful, or discriminatory.
- Making changes to production-critical files without proper testing.
- Adding features outside the scope of `users`, `transfers`, and `point_ledger` without approval.

## Additional Notes
- Always validate changes in a staging environment before deploying to production.
- Follow the repository's contribution guidelines and code review process.
- Ensure compatibility with the existing tech stack (Go, SQLite, Fiber).