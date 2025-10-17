# KBTG AI Backend Workshop

## Overview
This project demonstrates a backend application built with Go and Fiber framework. It includes CRUD operations for managing user data stored in an SQLite database. The user data fields align with the provided UI design.

## Features
- **GET /users**: Retrieve all users.
- **GET /users/{id}**: Retrieve details of a specific user.
- **POST /users**: Create a new user.
- **PUT /users/{id}**: Update an existing user's information.
- **DELETE /users/{id}**: Delete a user.

## Database Schema
The `users` table contains the following fields:
- `id`: Integer (Primary Key, Auto Increment)
- `membership`: Text
- `name`: Text
- `surname`: Text
- `phone`: Text
- `email`: Text
- `join_date`: Text
- `membership_level`: Text
- `points`: Integer

## Setup Instructions
1. Clone the repository:
   ```bash
   git clone https://github.com/NaponTunglukmongkol/kbtg-ai-backend-workshop.git
   ```

2. Navigate to the project directory:
   ```bash
   cd kbtg-ai-backend-workshop
   ```

3. Install dependencies:
   ```bash
   go mod tidy
   ```

4. Run the application:
   ```bash
   go run main.go
   ```

## Testing Endpoints
Use `curl` or any API testing tool (e.g., Postman) to test the endpoints:

### Example Requests
- **GET /users**:
  ```bash
  curl http://127.0.0.1:3000/users
  ```

- **POST /users**:
  ```bash
  curl -X POST http://127.0.0.1:3000/users \
       -H "Content-Type: application/json" \
       -d '{"membership":"Gold","name":"สมชาย","surname":"ใจดี","phone":"081-234-5678","email":"somchai@example.com","join_date":"2025-10-17","membership_level":"Gold","points":15420}'
  ```

- **GET /users/{id}**:
  ```bash
  curl http://127.0.0.1:3000/users/1
  ```

- **PUT /users/{id}**:
  ```bash
  curl -X PUT http://127.0.0.1:3000/users/1 \
       -H "Content-Type: application/json" \
       -d '{"membership":"Platinum","name":"สมชาย","surname":"ใจดี","phone":"081-234-5678","email":"somchai@example.com","join_date":"2025-10-17","membership_level":"Platinum","points":20000}'
  ```

- **DELETE /users/{id}**:
  ```bash
  curl -X DELETE http://127.0.0.1:3000/users/1
  ```

## License
This project is licensed under the MIT License.