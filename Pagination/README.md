# Product API with Pagination, Sorting, and Filtering

A Go-based REST API demonstrating advanced database querying techniques including pagination, sorting, and filtering. Built with Echo framework and GORM ORM for PostgreSQL.

## Learning Objectives

This project illustrates key concepts in building scalable APIs:

- **Pagination**: Implementing offset-based pagination for large datasets
- **Sorting**: Dynamic sorting by any field in ascending/descending order
- **Filtering**: Text-based search across multiple fields
- **Query Parameters**: Handling and validating URL query parameters
- **Database Optimization**: Efficient querying with LIMIT, OFFSET, and WHERE clauses
- **API Design**: RESTful endpoints with structured JSON responses

## Features

- **Product Management**: CRUD operations for product data
- **Data Seeding**: Populate database with sample products
- **Advanced Querying**: Paginated, sorted, and filtered product listings
- **Flexible Sorting**: Sort by any product field (name, category, price, etc.)
- **Text Filtering**: Search products by name or category
- **Pagination Metadata**: Returns page info, total items, and total pages

## Prerequisites

- Go 1.25.0 or later
- PostgreSQL database
- Git (for cloning the repository)

## Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd pagination
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
   SERVER_PORT=8080
   ```

4. **Set up PostgreSQL database:**

   Create a PostgreSQL database and update the `.env` file with your credentials.

## Database Schema

The application uses a `products` table with the following structure:

```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(255),
    price DECIMAL(10,2),
    description TEXT
);
```

GORM automatically creates this table via auto-migration.

## API Endpoints

### Seed Sample Data
- **Endpoint:** `POST /save`
- **Description:** Populates the database with sample product data
- **Response:**
  ```json
  {
    "message": "Data seeded successfully!"
  }
  ```

### Get Products (with Pagination, Sorting, Filtering)
- **Endpoint:** `GET /products`
- **Query Parameters:**
  - `page` (optional): Page number (default: 1)
  - `limit` (optional): Items per page (default: 5)
  - `sortField` (optional): Field to sort by (e.g., "name", "price", "category")
  - `sortOrder` (optional): "asc" or "desc" (default: "asc")
  - `filter` (optional): Search text for name or category filtering
- **Response:**
  ```json
  {
    "page": 1,
    "limit": 5,
    "total_items": 25,
    "total_pages": 5,
    "data": [
      {
        "id": 1,
        "name": "Laptop",
        "category": "Electronics",
        "price": 999.99,
        "description": "A high-performance laptop."
      }
    ]
  }
  ```

## Usage Examples

1. **Start the server:**
   ```bash
   go run main.go
   ```

2. **Seed data:**
   ```bash
   curl -X POST http://localhost:8080/save
   ```

3. **Get first page of products:**
   ```bash
   curl "http://localhost:8080/products?page=1&limit=5"
   ```

4. **Get products sorted by price (descending):**
   ```bash
   curl "http://localhost:8080/products?sortField=price&sortOrder=desc"
   ```

5. **Filter products by category:**
   ```bash
   curl "http://localhost:8080/products?filter=electronics"
   ```

6. **Combined pagination, sorting, and filtering:**
   ```bash
   curl "http://localhost:8080/products?page=2&limit=3&sortField=name&sortOrder=asc&filter=home"
   ```

## Code Structure

```
pagination/
├── main.go          # Main application with API endpoints
├── go.mod           # Go module dependencies
└── .env             # Environment variables (create this)
```

## Key Learning Concepts

### 1. Pagination Implementation
- **Offset Calculation**: `offset = (page - 1) * limit`
- **Total Pages**: `(total + limit - 1) / limit` (ceiling division)
- **GORM Methods**: `Limit()`, `Offset()`, `Count()`

### 2. Dynamic Sorting
- **Safe String Formatting**: Using `fmt.Sprintf()` for ORDER BY clauses
- **Order Validation**: Converting to lowercase and defaulting to "asc"
- **Field Flexibility**: Allowing any model field for sorting

### 3. Text Filtering
- **LIKE Queries**: Using `%` wildcards for partial matches
- **Case Insensitive**: Converting filter to lowercase with `strings.ToLower()`
- **Multiple Fields**: Searching across name and category columns

### 4. Query Parameter Handling
- **String to Integer**: Using `strconv.Atoi()` with error handling
- **Default Values**: Providing sensible defaults for missing parameters
- **Validation**: Ensuring positive values for page and limit

### 5. GORM Query Building
- **Method Chaining**: Building complex queries fluently
- **Conditional Queries**: Adding WHERE and ORDER clauses based on parameters
- **Count Before Limit**: Getting total count before applying pagination

## Dependencies

- **Echo v4**: Web framework for building HTTP services
- **GORM**: Powerful ORM library for Go
- **PostgreSQL Driver**: Database driver for PostgreSQL connections
- **godotenv**: Loads environment variables from .env files

## Performance Considerations

- **Database Indexing**: Consider adding indexes on frequently sorted/filtered columns
- **Query Optimization**: For very large datasets, consider cursor-based pagination
- **Caching**: Implement caching for frequently accessed pages
- **Rate Limiting**: Add rate limiting to prevent abuse

## Testing the API

You can test the endpoints using tools like:
- **curl**: Command-line HTTP client
- **Postman**: GUI API testing tool
- **Thunder Client**: VS Code extension for API testing

## Extending the Project

Ideas for enhancement:
- Add more filter options (price range, date filters)
- Implement cursor-based pagination for better performance
- Add authentication and authorization
- Create frontend interface for the API
- Add unit tests and integration tests
- Implement caching with Redis
- Add API documentation with Swagger

## Common Issues

- **Port already in use**: Change SERVER_PORT in .env or kill the process using the port
- **Database connection failed**: Verify PostgreSQL is running and credentials are correct
- **No data returned**: Make sure to call the /save endpoint first to seed data

## License

This project is for educational purposes. Feel free to modify and learn from the code.

---

## 🚀 How to Test

1.  **Seed Data:** Send a `POST` request to `/save` to populate the database.
2.  **Fetch with Logic:** Send a `GET` request to `/products` with parameters:
    * `GET /products?page=1&limit=2&filter=laptop&sortField=price&sortOrder=desc`



---

## 🏗 Future Enhancements
* **Validation:** Use a library like `go-playground/validator` to ensure product data is correct before saving.
* **Indexes:** Add indexes to the `Name` and `Category` columns in the database to speed up `LIKE` queries as the dataset grows.
* **Logging:** Implement structured logging to track API performance and errors in production.
