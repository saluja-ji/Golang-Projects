# File Upload Service in Go

A simple file upload and retrieval service built with Go, Echo framework, and PostgreSQL. This project demonstrates fundamental concepts in web development, database integration, and file handling.

## Learning Objectives

This project covers the following Go and web development concepts:
- Building REST APIs with Echo framework
- Database integration using GORM ORM
- File upload and download handling
- Environment variable configuration
- Database migrations with GORM
- HTTP request/response handling
- Error handling in web applications

## Features

- **File Upload**: Upload files via HTTP POST request
- **File Retrieval**: Download files by ID via HTTP GET request
- **Database Storage**: Files stored in PostgreSQL database
- **Metadata Tracking**: Store filename and content type
- **RESTful API**: Clean API endpoints for file operations

## Prerequisites

- Go 1.25.0 or later
- PostgreSQL database
- Git (for cloning the repository)

## Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd fileupload
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Set up environment variables:**

   Create a `.env` file in the root directory:
   ```env
   DB_HOST=localhost
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=your_db_name
   DB_PORT=5432
   ```

4. **Set up PostgreSQL database:**

   Create a PostgreSQL database and update the `.env` file with your database credentials.

## Database Schema

The application uses a single table `file_uploads` with the following structure:

```sql
CREATE TABLE file_uploads (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    file_name VARCHAR(255),
    file_type VARCHAR(255),
    file_data BYTEA
);
```

GORM automatically creates this table when the application starts.

## API Endpoints

### Upload File
- **Endpoint:** `POST /upload`
- **Content-Type:** `multipart/form-data`
- **Parameters:**
  - `file` (required): The file to upload
- **Response:**
  ```json
  {
    "message": "File saved successfully",
    "id": 1,
    "file_name": "example.txt",
    "file_type": "text/plain"
  }
  ```

### Retrieve File
- **Endpoint:** `GET /file/:id`
- **Parameters:**
  - `id` (path parameter): The file ID
- **Response:** The file content with appropriate content-type header

## Usage

1. **Start the server:**
   ```bash
   go run main.go
   ```

2. **Upload a file:**
   ```bash
   curl -X POST -F "file=@/path/to/your/file.txt" http://localhost:8082/upload
   ```

3. **Download a file:**
   ```bash
   curl http://localhost:8082/file/1 -o downloaded_file.txt
   ```

## Code Structure

```
fileupload/
├── main.go          # Main application file
├── go.mod           # Go module dependencies
└── .env             # Environment variables (create this)
```

## Key Learning Points

### 1. Echo Framework
- Setting up an Echo server
- Defining routes and handlers
- Using middleware for logging and recovery

### 2. File Handling
- Reading multipart form data
- Converting file to byte array
- Streaming file content in responses

### 3. Database Operations
- Connecting to PostgreSQL with GORM
- Auto-migration of structs to tables
- CRUD operations with GORM

### 4. Environment Configuration
- Loading environment variables with godotenv
- Secure configuration management

### 5. Error Handling
- Proper HTTP status codes
- JSON error responses
- Graceful error handling

## Dependencies

- **Echo v4**: Web framework for Go
- **GORM**: ORM library for Go
- **PostgreSQL Driver**: Database driver for PostgreSQL
- **godotenv**: Environment variable loader

## Running Tests

Currently, this project doesn't include unit tests. As a learning exercise, you could add tests for:
- File upload functionality
- File retrieval functionality
- Database operations
- Error handling scenarios

## Contributing

This is a learning project. Feel free to:
- Add new features
- Improve error handling
- Add tests
- Optimize performance
- Add more file type validations

## License

This project is for educational purposes only.
