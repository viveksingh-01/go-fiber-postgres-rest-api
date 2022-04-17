package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/viveksingh-01/go-fiber-postgres-rest-api/database"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) SetupRoutes(app *fiber.App) {}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &database.Config{
		Host: os.Getenv("DB_HOST"),
		Port: os.Getenv("DB_PORT"),
		DBName: os.Getenv("DB_NAME"),
		User: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	}
	db, err := database.Connect(config)

	r := Repository {
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Get("/healthcheck", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).SendString("All good here!")
	})
	log.Fatal(app.Listen(":5000"))
}
