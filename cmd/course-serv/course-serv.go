package main

import (
	"log"
	"os"

	"github.com/estradax/course-serv/internal"
	"github.com/estradax/course-serv/internal/handler"
	"github.com/gofiber/fiber/v2"
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

	registerHandler := handler.NewRegisterHandler(db, []byte(jwtSecret))

	app := fiber.New()

	app.Post("/api/v1/register", registerHandler.Register)

	app.Post("/api/v1/login", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world")
	})

	if err := app.Listen(":8080"); err != nil {
		log.Fatalln("Cannot listen: ", err.Error())
	}
}
