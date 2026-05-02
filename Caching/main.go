package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/redis/go-redis/v9"
)

var (
	db  *gorm.DB
	rds *redis.Client
	ctx = context.Background()
)

type User struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	NAME  string `json:"name"`
	EMAIL string `json:"email" gorm:"unique"`
}

func initDB() *gorm.DB {
	dsn := os.Getenv("POSTGRES_DNS")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}
	db.AutoMigrate(&User{})
	return db
}

func initRedis() *redis.Client {
	// options here help us in keeping the seeting to connect to the redis server
	// jsut like we do in postgres, we can set the address, password and db number to connect to the redis server
	rds := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	// result of ping command is used to check if the connection to the redis server is successful or not
	_, err := rds.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	return rds
}

func getAllUsers(c echo.Context) error {
	// here get is used to get the value of the key "users" from the redis server,
	//  if the key is not present in the redis server, then it will return an error
	val, err := rds.Get(ctx, "all_users").Result()
	if err == redis.Nil {
		var users []User
		if err := db.Find(&users).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch users from database!"})
		}
		// here marshal is used to convert the users slice into a json format, so that it can be stored in the redis server as a string
		data, _ := json.Marshal(users)
		// here set is used to set the value of the key "users" in the redis server, with the json data of users and an expiration time of 1 hour
		//Set Redis `SET key value [expiration]` command. Use expiration for `SETEx`-like behavior.
		rds.Set(ctx, "all_users", data, 10*time.Minute)

		return c.JSON(http.StatusOK, users)
	} else if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch users from cache!"})
	}

	var users []User
	// here unmarshal is used to convert the json data of users from the redis server back into a slice of User struct
	if err := json.Unmarshal([]byte(val), &users); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to parse users from cache!"})
	}
	return c.JSON(http.StatusOK, users)
}

func createUser(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request"})
	}

	var lastUser User

	if err := db.Order("id desc").First(&lastUser).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"erro": "DB errro"})
	}
	u.ID = lastUser.ID + 1
	if err := db.Create(u).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create user in database!"})
	}

	var users []User
	if err := db.Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to refresh user cache!"})
	}
	data, _ := json.Marshal(users)
	rds.Set(ctx, "all_users", data, 10*time.Minute)

	data, _ = json.Marshal(u)
	rds.Set(ctx, fmt.Sprintf("user:%d", u.ID), data, 10*time.Minute)
	return c.JSON(http.StatusCreated, u)
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}
	db = initDB()
	rds = initRedis()

	e := echo.New()
	e.GET("/allUsers", getAllUsers)
	e.POST("/user", createUser)

	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))
}
