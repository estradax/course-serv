package main

import (
	"log"

	"github.com/estradax/course-serv/internal"
	"github.com/gofiber/fiber/v2"
)

func main() {
	_, err := internal.ConnectDB()
	if err != nil {
		log.Fatalln("Something went wrong: ", err.Error())
	}

	app := fiber.New()

	app.Post("/api/v1/register", func(c *fiber.Ctx) error {
		return c.SendString("Register")
	})

	app.Post("/api/v1/login", func(c *fiber.Ctx) error {
		return c.SendString("Hello, world")
	})

	app.Listen(":8080")
}
