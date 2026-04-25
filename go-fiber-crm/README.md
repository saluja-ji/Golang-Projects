# Go Fiber CRM Project

## What is this project?

This is a simple **CRM (Customer Relationship Management)** application built using **Go** programming language. A CRM is a tool that helps businesses manage their interactions with current and potential customers. In this project, we're managing "leads" – which are basically potential customers or contacts.

Think of it like a digital address book for business contacts, but with the ability to add, view, and remove entries through a web API.

## Prerequisites (What you need before starting)

Before you can run this project, you need to have **Go** installed on your computer. Go is a programming language created by Google.

### How to install Go

1. Go to the official Go website: https://golang.org/dl/
2. Download the installer for your operating system (Windows, macOS, or Linux)
3. Run the installer and follow the instructions
4. Open a terminal/command prompt and type `go version` to check if it's installed correctly

## Project Structure

Let's break down what each file does:

```
go-fiber-crm/
├── go.mod          # This file tells Go what libraries our project needs
├── main.go         # The main entry point of our application
├── database/
│   └── database.go # Handles database connection
└── lead/
    └── lead.go     # Contains the Lead model and API handlers
```

### What is go.mod?

The `go.mod` file is like a shopping list for our project. It tells Go:
- What version of Go we're using (1.25.0)
- What external libraries we need (like Fiber and GORM)

### What is main.go?

This is the "brain" of our application. It:
- Sets up the web server using Fiber
- Connects to the database
- Defines the API routes (URLs that our app responds to)
- Starts the server on port 3000

### What is database/database.go?

This file is very simple. It just declares a global variable `DBcom` that will hold our database connection. Think of it as a phone line that connects our app to the database.

### What is lead/lead.go?

This file defines:
- The `Lead` struct (a blueprint for what a lead looks like)
- Functions to handle API requests (get leads, create lead, etc.)

## What are the technologies used?

### Go (Golang)
Go is a modern programming language that's fast, reliable, and easy to learn. It's great for building web servers and APIs.

### Fiber
Fiber is a web framework for Go. Frameworks provide ready-made tools so you don't have to build everything from scratch. Fiber helps us create web routes and handle HTTP requests easily.

### GORM
GORM is an ORM (Object-Relational Mapping) library. It makes it easy to work with databases by letting us use Go structs instead of writing raw SQL queries.

### SQLite
SQLite is a simple database that stores data in a single file. It's perfect for small projects and doesn't require a separate database server.

## How to run the project

1. **Clone or download this project** to your computer

2. **Open a terminal** and navigate to the project folder:
   ```
   cd go-fiber-crm
   ```

3. **Download the dependencies** (libraries listed in go.mod):
   ```
   go mod download
   ```

4. **Run the application**:
   ```
   go run main.go
   ```

5. The server will start and you'll see messages like:
   ```
   connection open to database
   Database migrated
   ```

6. Your API is now running at `http://localhost:3000`

## API Endpoints (How to use the CRM)

Once the server is running, you can interact with it using tools like:
- **curl** (command line)
- **Postman** (GUI application)
- **Browser** (for GET requests)
- Or any HTTP client

### Get all leads
```
GET http://localhost:3000/api/v1/lead
```
This returns a list of all leads in JSON format.

### Get a specific lead
```
GET http://localhost:3000/api/v1/lead/1
```
Replace `1` with the ID of the lead you want to see.

### Create a new lead
```
POST http://localhost:3000/api/v1/lead
Content-Type: application/json

{
  "name": "John Doe",
  "company": "ABC Corp",
  "email": "john@example.com",
  "phone": 1234567890
}
```

### Delete a lead
```
DELETE http://localhost:3000/api/v1/lead/1
```
Replace `1` with the ID of the lead you want to delete.

## Code Walkthrough

### The Lead struct
In `lead/lead.go`, we define what a Lead looks like:

```go
type Lead struct {
    gorm.Model        // This adds ID, CreatedAt, UpdatedAt, DeletedAt fields automatically
    Name    string `json:"name"`     // The person's name
    Company string `json:"company"`  // Their company
    Email   string `json:"email"`    // Email address
    Phone   int    `json:"phone"`    // Phone number
}
```

The `json:"name"` tags tell Go how to convert this struct to/from JSON when sending data over the web.

### Database Connection
In `main.go`, we connect to SQLite:

```go
database.DBcom, err = gorm.Open(sqlite.Open("leads.db"), &gorm.Config{})
```

This creates a file called `leads.db` in your project folder if it doesn't exist.

### Auto Migration
```go
database.DBcom.AutoMigrate(&lead.Lead{})
```

This automatically creates the database table based on our Lead struct. It's like telling the database "create a table that matches this Go struct".

### API Handlers
Each API endpoint has a corresponding function:

- `GetLeads()`: Gets all leads from database and returns as JSON
- `GetLead()`: Gets one lead by ID
- `NewLead()`: Creates a new lead from JSON data in the request body
- `DeleteLead()`: Deletes a lead by ID

## What happens when you run the app?

1. **Database setup**: Connects to SQLite and creates the `leads` table
2. **Routes setup**: Tells Fiber which functions to call for each URL
3. **Server starts**: Listens for HTTP requests on port 3000
4. **Ready to use**: You can now send requests to the API endpoints

## Common issues and solutions

### "go: command not found"
You need to install Go first. Follow the installation steps above.

### "module not found" errors
Run `go mod download` to download dependencies.

### Port 3000 already in use
Change the port in `main.go` from `:3000` to something else like `:3001`.

### Database file not created
Check that you have write permissions in the project folder.

## Next steps

This is a basic CRM. You could extend it by:
- Adding user authentication
- Adding more fields to the Lead struct
- Creating a web interface (HTML/CSS/JavaScript)
- Adding search functionality
- Using a different database (PostgreSQL, MySQL)

## Learning resources

- [Go official tutorial](https://tour.golang.org/)
- [Fiber documentation](https://docs.gofiber.io/)
- [GORM documentation](https://gorm.io/)
- [REST API basics](https://restfulapi.net/)

Happy coding! 🚀