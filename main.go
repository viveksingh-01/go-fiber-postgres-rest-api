package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/viveksingh-01/go-fiber-postgres-rest-api/database"
	"github.com/viveksingh-01/go-fiber-postgres-rest-api/models"
	"gorm.io/gorm"
)

type Repository struct {
	DB *gorm.DB
}

type Book struct {
	Title     string `json:"title"`
	Author    string `json:"author"`
	Publisher string `json:"publisher"`
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/healthcheck", checkAPIHealth)
	api.Post("/books", r.AddBook)
	api.Get("/books", r.GetBooks)
}

func (r *Repository) AddBook(c *fiber.Ctx) error {
	book := Book{}
	err := c.BodyParser(&book)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "Request failed."})
		return err
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Could not add book."})
		return err
	}
	c.Status(http.StatusCreated).JSON(&fiber.Map{"message": "Book was added successfully!"})
	return nil
}

func (r *Repository) GetBooks(c *fiber.Ctx) error {
	books := &[]models.Books{}
	err := r.DB.Find(books).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "Couldn't fetch the books from the database."})
		return err
	} else {
		c.Status(http.StatusOK).JSON(&fiber.Map{
			"message": "Books were fetched successfully!",
			"data": books,
		})
	}
	return nil
}

func checkAPIHealth(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).SendString("All good here!")
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBName:   os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db, err := database.Connect(config)
	if err != nil {
		log.Fatal("Couldn't connect to the database.")
	} else {
		fmt.Println("Connected to the database successfully!")
	}

	err = models.MigrateBooks(db)
	if err != nil {
		log.Fatal("Couldn't migrate Books to the database.")
	}

	r := Repository {
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	log.Fatal(app.Listen(":5000"))
}
