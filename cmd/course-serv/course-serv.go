package main

import (
	"log"
	"os"

	"github.com/estradax/course-serv/internal"
	"github.com/estradax/course-serv/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Cannot load environment: ", err.Error())
	}

	db, err := internal.ConnectDB()
	if err != nil {
		log.Fatalln("Something went wrong: ", err.Error())
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	authHandler := handler.NewAuthHandler(db, []byte(jwtSecret))

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/v1/profile", authHandler.Profile)

	app.Post("/api/v1/register", authHandler.Register)

	app.Post("/api/v1/login", authHandler.Login)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalln("Cannot listen: ", err.Error())
	}
}
