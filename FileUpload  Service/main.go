package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type FileUpload struct {
	gorm.Model
	FileName string `json:"file_name,omitempty"`
	FileType string `json:"file_type,omitempty"`
	FileData []byte `json:"-" gorm:type:bytes"`
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

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&FileUpload{})
	db = database
}

func UploadFile(c echo.Context) error { // this is a Echo handler function c echo.Context gives us access to the request and response objects, as well as other context information about the request.
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "File is required"})
	}
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Unable to open file"})
	}
	defer src.Close()

	fileBytes, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Unable to read file"})
	}

	fileUpload := FileUpload{
		FileName: file.Filename,
		FileType: file.Header.Get("Content-Type"),
		FileData: fileBytes,
	}

	if err := db.Create(&fileUpload).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Unable to save file"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":   "File saved successfully",
		"id":        fileUpload.ID,
		"file_name": fileUpload.FileName,
		"file_type": fileUpload.FileType,
	})
}

func GetFile(c echo.Context) error {
	id := c.Param("id")
	var file FileUpload

	if err := db.First(&file, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "File not found"})
	}

	return c.Stream(http.StatusOK, file.FileType, bytes.NewReader(file.FileData))
}

func main() {
	initDB()
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.POST("/upload", UploadFile)
	e.GET("/file/:id", GetFile)

	e.Logger.Fatal(e.Start(":8082"))
}
