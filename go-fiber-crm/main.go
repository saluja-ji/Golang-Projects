package main

import (
	"fmt"
	"go-fiber-crm/database"
	"go-fiber-crm/lead"

	"github.com/gofiber/fiber"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRoutes(app *fiber.App) {
	app.Get("/api/v1/lead", lead.GetLeads)
	app.Get("/api/v1/lead/:id", lead.GetLead)
	app.Post("/api/v1/lead", lead.NewLead)
	app.Delete("/api/v1/lead/:id", lead.DeleteLead)
}

func initDatabase() {
	var err error

	database.DBcom, err = gorm.Open(sqlite.Open("leads.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("connection open to database")

	database.DBcom.AutoMigrate(&lead.Lead{})
	fmt.Println("Database migrated")
}

func main() {
	app := fiber.New()

	initDatabase()
	setupRoutes(app)

	sqlDB, _ := database.DBcom.DB()
	defer sqlDB.Close()

	app.Listen(":3000")
}
