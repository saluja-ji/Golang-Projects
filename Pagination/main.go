package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	Category    string  `json:"category"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

var db *gorm.DB

func initDB() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable ",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	fmt.Println("Connecting with DSN:", dsn)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Product{})
	db = database
}

// echo.Context is a toolbox for the current HTTP request that lets you read input and send output easily.
// It contains methods for request and response handling, path parameters, query parameters, and more.
func SeedData(c echo.Context) error {
	sample := []Product{
		{Name: "Laptop", Category: "Electronics", Price: 999.99, Description: "A high-performance laptop."},
		{Name: "Smartphone", Category: "Electronics", Price: 499.99, Description: "A powerful smartphone."},
		{Name: "Headphones", Category: "Electronics", Price: 199.99, Description: "Noise-cancelling headphones."},
		{Name: "Coffee Maker", Category: "Home Appliances", Price: 79.99, Description: "Brew the perfect cup of coffee."},
		{Name: "Blender", Category: "Home Appliances", Price: 59.99, Description: "Blend smoothies and more."},
	}
	db.Create(&sample)
	return c.JSON(http.StatusOK, echo.Map{"message": "Data seeded successfully!"})
}

func GetProducts(c echo.Context) error {
	pageParam := c.QueryParam("page") // over here we are  extracting pagination inputs (page, limit) from the API request URL.
	limitParams := c.QueryParam("limit")
	sortFiled := c.QueryParam("sortField")
	sortOrder := c.QueryParam("sortOrder")
	filter := c.QueryParam("filter")

	// first we convet page into interger
	// strconv.Atoi() converts a string into an integer (int).
	page, err := strconv.Atoi(pageParam)
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(limitParams)
	if err != nil || limit <= 0 {
		limit = 5
	}

	offset := (page - 1) * limit

	query := db.Model(&Product{})
	if filter != "" {
		// % means match any sequence of characters
		filterPattern := "%" + strings.ToLower(filter) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(category) LIKE ?", filterPattern, filterPattern)
	}

	if sortFiled != "" {
		order := "asc"
		if strings.ToLower(sortOrder) == "desc" {
			order = "desc"
		}
		query = query.Order(fmt.Sprintf("%s %s", sortFiled, order))
	}
	var total int64
	query.Count(&total)

	var product []Product
	if err := query.Limit(limit).Offset(offset).Find(&product).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch the bulshit data "})
	}

	totalPages := (int(total) + limit - 1) / limit

	return c.JSON(http.StatusOK, echo.Map{
		"page":        page,
		"limit":       limit,
		"total_items": total,
		"total_pages": totalPages,
		"data":        product,
	})
}

func main() {
	initDB()

	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.POST("/save", SeedData)

	e.GET("/products", GetProducts)

	port := os.Getenv("SERVER_PORT")

	e.Logger.Fatal(e.Start(":" + port))
}
