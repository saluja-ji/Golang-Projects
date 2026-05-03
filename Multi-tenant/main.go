package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"          // JWT creation + validation
	echojwt "github.com/labstack/echo-jwt/v4" // Middleware for JWT auth in Echo
	"github.com/labstack/echo/v4"           // Web framework
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"              // PostgreSQL driver
	"gorm.io/gorm"                         // ORM for DB operations
)

// Represents a user table in DB
type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"` // auto-increment primary key
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"` // ⚠️ BAD: storing plain text passwords
	Role      string `json:"role"`     // used for authorization (admin/tenant/user)
	CreatedAt string `json:"created_at"`
}

// Custom JWT payload (claims)
// Combines your custom data + standard JWT fields
type JwtCustomClaim struct {
	UserID uint   `json:"user_id"` // identifies user making request
	Role   string `json:"role"`    // used later for access control

	// Embedded struct → gives access to standard claims like exp, iat
	jwt.RegisteredClaims
}

var (
	db         *gorm.DB
	jwtSecrets = "fjkenrjifkje3k45k34tk" // ⚠️ BAD: hardcoded secret (should be env variable)
)

// Initialize database connection
func initDB() {
	dsn := "host=localhost user=xxxx password=xxxx dbname=xxx port=5432 sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!") // crash app if DB fails
	}

	// Auto-create/update schema
	db.AutoMigrate(&User{})
}

// ----------------------
// REGISTER USER
// ----------------------
func register(c echo.Context) error {

	// Struct for incoming JSON request
	type Input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	var input Input

	// Bind request body → maps JSON → Go struct
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Create user object
	user := User{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password, // ⚠️ SHOULD HASH (bcrypt)
		Role:     input.Role,
	}

	// Insert into DB
	if err := db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// ----------------------
// LOGIN → GENERATE JWT
// ----------------------
func login(c echo.Context) error {

	type Input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input Input

	// Parse request JSON
	if err := c.Bind(&input); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var u User

	// Find user by email
	if err := db.Where("email=?", input.Email).First(&u).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	// ⚠️ Plain text comparison → insecure
	if u.Password != input.Password {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "password does not match"})
	}

	// Build JWT payload
	claims := JwtCustomClaim{
		UserID: u.ID,
		Role:   u.Role,

		// Standard claims (important for security)
		RegisteredClaims: jwt.RegisteredClaims{

			// Token expires in 24 hours → prevents reuse forever
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),

			// Timestamp when token created
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token object with HS256 signing algorithm
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key → produces final JWT string
	t, err := token.SignedString([]byte(jwtSecrets))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"Token": t})
}

// ----------------------
// ADMIN ROUTE
// ----------------------
func adminDashboard(c echo.Context) error {

	// Extract token from context (added by middleware)
	user := c.Get("user").(*jwt.Token)

	// Cast claims back to your struct
	claims := user.Claims.(*JwtCustomClaim)

	// Authorization check
	if claims.Role != "admin" {
		return echo.NewHTTPError(http.StatusForbidden, "only for admin access")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"messagee": "Welcome to the admin dashboard!",
		"user_id":  claims.UserID,
	})
}

// ----------------------
// TENANT ROUTE
// ----------------------
func tenantDashboard(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaim)

	if claims.Role != "tenant" {
		return echo.NewHTTPError(http.StatusForbidden, "only for tenant access")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"messagee": "Welcome to the tenant dashboard!",
		"user_id":  claims.UserID,
	})
}

// ----------------------
// USER ROUTE
// ----------------------
func userDashboard(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaim)

	if claims.Role != "user" {
		return echo.NewHTTPError(http.StatusForbidden, "only for user access")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"messagee": "Welcome to the user dashboard!",
		"user_id":  claims.UserID,
	})
}

// ----------------------
// MAIN FUNCTION
// ----------------------
func main() {

	initDB()

	e := echo.New()

	// Logs requests (method, path, status)
	e.Use(middleware.RequestLogger())

	// Recovers from panics → prevents server crash
	e.Use(middleware.Recover())

	// Public routes
	e.POST("/register", register)
	e.POST("/login", login)

	// JWT middleware config
	config := echojwt.Config{

		// Secret used to verify token signature
		SigningKey: []byte(jwtSecrets),

		// Tell middleware what struct to decode claims into
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(JwtCustomClaim)
		},
	}

	// Protected route group
	r := e.Group("/api")

	// Apply JWT middleware → ALL routes below require valid token
	r.Use(echojwt.WithConfig(config))

	// Role-based endpoints
	r.POST("/admin", adminDashboard)
	r.POST("/tenant", tenantDashboard)
	r.POST("/user", userDashboard)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}
