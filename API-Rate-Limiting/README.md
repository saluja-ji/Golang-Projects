# API Rate Limiting Project

This is a simple Go project demonstrating API rate limiting using the Echo web framework and GORM for database interactions. It's designed for learning purposes to understand how to implement rate limiting in a REST API.

## What This Project Does

This project creates a basic API server with:
- User registration endpoint (`/signup`) that generates unique API keys
- A protected data endpoint (`/data`) that is rate-limited
- Rate limiting: 15 requests per 15-second window per API key
- PostgreSQL database to store users and rate limit counters

## Learning Objectives

By studying this code, you'll learn:
- How to build REST APIs with Echo framework
- Database integration using GORM
- Implementing middleware for cross-cutting concerns like rate limiting
- API key authentication
- Sliding window rate limiting algorithm
- HTTP headers and status codes

## Prerequisites

- Go 1.19 or later
- PostgreSQL database
- Basic understanding of Go, HTTP, and databases

## Setup

1. **Clone or navigate to the project directory:**
   ```bash
   cd API_rate_limiting
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up PostgreSQL database:**
   - Create a database named `mydb`
   - Update the database connection string in `main.go` 
4. **Run the server:**
   ```bash
   go run main.go
   ```

The server will start on port 8080.

## Usage

### Register a User
```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe"}'
```

Response:
```json
{
  "message": "user created successfully",
  "api_key": "API-6e22a529088475d1f2a04c11488cc8cc"
}
```

### Access Protected Data
```bash
curl -X GET http://localhost:8080/data \
  -H "X-API-Key: API-6e22a529088475d1f2a04c11488cc8cc"
```

Response (within limit):
```json
{
  "message": "welcome to the data api",
  "time": "2023-10-01T12:00:00Z"
}
```

Response (rate limit exceeded):
```json
{
  "error": "Rate limit exceeded",
  "retry_after": 10,
  "status": "exceeded"
}
```

## How Rate Limiting Works

The `RateLimiter` function creates middleware that:
1. Checks for an API key in the `X-API-Key` header
2. Tracks request counts per API key in a sliding window
3. Blocks requests when the limit (15 per 15 seconds) is exceeded
4. Returns appropriate HTTP status codes and retry information

## Code Structure

- `main.go`: Contains all the application logic
- `User` struct: Represents API users
- `RateLimit` struct: Stores rate limiting data
- `RateLimiter()`: Factory function that creates rate limiting middleware

## Next Steps for Learning

- Experiment with different rate limits and window sizes
- Add more endpoints with different rate limits
- Implement different rate limiting algorithms (token bucket, leaky bucket)
- Add logging and monitoring
- Create unit tests for the rate limiting logic

## Dependencies

- [Echo](https://echo.labstack.com/): Web framework
- [GORM](https://gorm.io/): ORM for Go
- [PostgreSQL driver](https://github.com/lib/pq): Database driver</content>
<filePath">/home/pushpit-saluja/go/src/API_rate_limiting/README.md
