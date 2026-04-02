# Go Movies CRUD

A simple REST API for managing a collection of movies, built with Go and the Gorilla Mux router. This project demonstrates basic CRUD operations (Create, Read, Update, Delete) for movie data stored in memory.

## Features

- Get all movies
- Get a specific movie by ID
- Create a new movie
- Update an existing movie
- Delete a movie
- JSON-based API responses

## Prerequisites

- Go 1.25.0 or later
- Gorilla Mux library (automatically installed via go.mod)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/go-movies-crud.git
   cd go-movies-crud
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Application

Start the server:
```bash
go run main.go
```

The server will start on port 8000. You can access the API at `http://localhost:8000`.

## API Endpoints

### Get All Movies
- **GET** `/movies`
- Returns a JSON array of all movies.

### Get Movie by ID
- **GET** `/movies/{id}`
- Returns a JSON object of the movie with the specified ID.

### Create Movie
- **POST** `/movies`
- Body: JSON object with movie details (isbn, title, director)
- Returns the created movie with a generated ID.

### Update Movie
- **PUT** `/movies/{id}`
- Body: JSON object with updated movie details
- Returns the updated movie.

### Delete Movie
- **DELETE** `/movies/{id}`
- Returns the updated list of movies after deletion.

## Sample Data

The application starts with two sample movies:
- "The great" by John Doe
- "The great economist" by George Soros

## Personal Notes

This project was created as a learning exercise to understand REST API development in Go. It uses in-memory storage, so data is lost when the server restarts. For production use, consider adding a database like PostgreSQL or MongoDB.

## Contributing

Feel free to fork and contribute! Open issues or submit pull requests for improvements.


