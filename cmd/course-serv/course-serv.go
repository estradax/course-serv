package main

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/estradax/course-serv/internal"
	"github.com/estradax/course-serv/internal/handler"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	err := internal.LoadEnv()
	if err != nil {
		log.Println("Cannot loadEnv: ", err.Error())
	}

	db, err := internal.ConnectDB()
	if err != nil {
		log.Fatalln("Something went wrong: ", err.Error())
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	cldUrl := os.Getenv("CLOUDINARY_URL")
	cld, err := cloudinary.NewFromURL(cldUrl)
	if err != nil {
		log.Fatalln("Cannot init cloudinary: ", err.Error())
	}

	h := handler.NewHandler(db, []byte(jwtSecret), cld)

	app := fiber.New()

	app.Use(cors.New())

	app.Get("/api/v1/profile", h.Authenticated, h.Profile)
	app.Post("/api/v1/register", h.Register)
	app.Post("/api/v1/login", h.Login)

	app.Get("/api/v1/courses", h.CourseGetAll)
	app.Get("/api/v1/recommended", h.CourseBasicToLearn)

	if err := app.Listen(":8080"); err != nil {
		log.Fatalln("Cannot listen: ", err.Error())
	}
}
