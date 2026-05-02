# Go Redis Caching Example

This is a simple Go web application that demonstrates how to integrate Redis caching with a PostgreSQL database using the Echo framework and GORM ORM. It's designed for learning purposes to understand caching strategies, database operations, and REST API development in Go.

## Features

- **User Management**: Create and retrieve users via REST API
- **Redis Caching**: Cache user data to reduce database load
- **PostgreSQL Database**: Persistent storage using GORM
- **RESTful API**: Clean endpoints for user operations
- **Environment Configuration**: Secure configuration using environment variables

## Technologies Used

- **Go**: Programming language
- **Echo**: Web framework for building REST APIs
- **GORM**: ORM for database operations
- **PostgreSQL**: Relational database
- **Redis**: In-memory data structure store for caching
- **godotenv**: Load environment variables from .env file

## Prerequisites

Before running this application, make sure you have the following installed:

- Go (version 1.19 or later)
- PostgreSQL database
- Redis server
- Git (for cloning the repository)

## Installation

1. **Clone the repository**:
   ```bash
   git clone <repository-url>
   cd REDIS
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up environment variables**:
   
   Create a `.env` file in the root directory with the following variables:
   ```
   POSTGRES_DNS=host=localhost user=your_username password=your_password dbname=your_db port=5432 sslmode=disable
   REDIS_ADDR=localhost:6379
   REDIS_PASSWORD=your_redis_password  # Leave empty if no password
   PORT=8080
   ```

4. **Set up PostgreSQL database**:
   
   Create a database in PostgreSQL and ensure the connection string in `.env` is correct.

5. **Start Redis server**:
   
   Make sure Redis is running on the specified address (default: localhost:6379).

## Running the Application

1. **Build and run**:
   ```bash
   go run main.go
   ```

   The server will start on the port specified in the `PORT` environment variable (default: 8080).

## API Endpoints

### GET /allUsers
Retrieves all users. First checks Redis cache, falls back to database if not cached.

**Response**: JSON array of users
```json
[
  {
    "ID": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
]
```

### POST /user
Creates a new user.

**Request Body**:
```json
{
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

**Response**: Created user object
```json
{
  "ID": 2,
  "name": "Jane Doe",
  "email": "jane@example.com"
}
```

## How Caching Works

This application demonstrates a simple caching strategy:

1. **Cache Miss**: When `/allUsers` is called and data isn't in Redis, it fetches from PostgreSQL, stores in Redis with a 10-minute expiration, then returns the data.

2. **Cache Hit**: Subsequent calls retrieve data directly from Redis, reducing database load.

3. **Cache Invalidation**: When a new user is created, the cache is updated with fresh data from the database.

## Learning Concepts

This project covers several important concepts in Go web development:

- **Database Integration**: Using GORM to interact with PostgreSQL
- **Caching Strategies**: Implementing Redis for performance optimization
- **REST API Design**: Building clean, RESTful endpoints
- **Environment Management**: Secure configuration with environment variables
- **Error Handling**: Proper error responses and logging
- **JSON Serialization**: Converting between Go structs and JSON
- **Middleware**: Using Echo's built-in features

## Project Structure

```
REDIS/
├── main.go          # Main application file
├── go.mod           # Go module file
├── go.sum           # Go dependencies
├── .env             # Environment variables (create this)
└── README.md        # This file
```

## Troubleshooting

- **Database Connection Issues**: Ensure PostgreSQL is running and credentials are correct
- **Redis Connection Issues**: Verify Redis server is running and address/password are correct
- **Port Already in Use**: Change the PORT in .env if 8080 is occupied
- **Module Issues**: Run `go mod tidy` to clean up dependencies

## Next Steps for Learning

To expand your knowledge, consider:

- Adding more endpoints (update, delete users)
- Implementing authentication/authorization
- Adding input validation
- Using Redis for more complex caching scenarios
- Adding logging middleware
- Writing unit tests
- Containerizing with Docker

## License

This project is for educational purposes. Feel free to modify and learn from it!
