# Multitenant JWT Authentication Example

This project is a simple learning example showing how to build a role-based authentication API in Go using:

- `Echo` framework for routing and HTTP handling
- `GORM` ORM for PostgreSQL database access
- `JWT` bearer tokens for authentication
- Role-based access control for protected endpoints

## What the code does

1. Connects to a PostgreSQL database using GORM.
2. Automatically creates a `users` table from the `User` struct model.
3. Offers two public endpoints:
   - `POST /register` to create a new user
   - `POST /login` to authenticate and receive a JWT token
4. Provides protected endpoints under `/api`:
   - `POST /api/admin` for `admin` role users only
   - `POST /api/tenant` for `tenant` role users only
   - `POST /api/user` for `user` role users only
5. Uses JWT claims to store the user ID and role, then validates that token on protected routes.

## Key structs and concepts

### `User`
The `User` struct holds basic user information:

- `ID` - auto-generated primary key
- `Name` - user name
- `Email` - login email
- `Password` - plain text password in this example (not secure)
- `Role` - role used for access control
- `CreatedAt` - optional timestamp metadata

### `JwtCustomClaim`
This struct extends the built-in JWT claims with:

- `UserID` - the authenticated user ID
- `Role` - the user role (`admin`, `tenant`, or `user`)

These custom claims are embedded inside every signed token.

## Important functions

### `initDB()`
- Builds the database connection string.
- Connects to PostgreSQL.
- Runs `AutoMigrate` so the `users` table exists.

### `register()`
- Reads user input from JSON body.
- Saves a new `User` record to the database.
- Returns the saved user as JSON.

> Note: Passwords are stored as plain text in this example for learning purposes only. In real applications, always hash passwords.

### `login()`
- Reads login credentials from JSON body.
- Finds the user by email and checks password.
- Creates a JWT token with `UserID`, `Role`, expiration, and issue time.
- Returns the signed token as JSON.

### `adminDashboard()`, `tenantDashboard()`, `userDashboard()`
- Extract the JWT token from the request context.
- Read the custom claims from the token.
- Check the role value and allow access only if the role matches.

## Middleware and routing

- Global middleware includes request logging and panic recovery.
- The `/api` route group is protected by `echo-jwt` middleware.
- Tokens are validated using the same `jwtSecrets` signing key.
- `NewClaimsFunc` ensures that JWT claims are parsed into `JwtCustomClaim` objects.

## How to run

1. Ensure PostgreSQL is running and accessible.
2. Update the DSN in `main.go` if needed:

```
host=localhost user=psssh password=p1223. dbname=mydb port=5432 sslmode=disable
```

3. Run the program:

```bash
go run main.go
```

4. Use a tool like `curl` or Postman to test the endpoints.

## Example flow

1. Register a user:

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret","role":"admin"}'
```

2. Login and get a token:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"alice@example.com","password":"secret"}'
```

3. Call a protected route:

```bash
curl -X POST http://localhost:8080/api/admin \
  -H "Authorization: Bearer <TOKEN>"
```

## Learning points

- How to wire a Go web server with Echo.
- How to define models and migrate a database using GORM.
- How JWT authentication works in a stateless API.
- How to use role-based authorization in protected endpoints.
- Why production apps need password hashing and stronger token management.
