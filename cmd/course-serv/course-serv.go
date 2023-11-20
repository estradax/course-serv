package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/estradax/course-serv/internal"
	"github.com/estradax/course-serv/internal/handler"
	"github.com/estradax/course-serv/internal/middleware"
	"github.com/estradax/course-serv/internal/model"
	"github.com/estradax/course-serv/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
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
	authService := service.NewAuthService(db, []byte(jwtSecret))
	courseService := service.NewCourseService(db, []byte(jwtSecret), cld)
	middlewareService := middleware.New(db, []byte(jwtSecret))

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(cors.New())

	app.Static("/", "./public")

	app.Get("/api/v1/profile", h.Authenticated, h.Profile)
	app.Post("/api/v1/register", h.Register)
	app.Post("/api/v1/login", h.Login)

	app.Get("/api/v1/courses", h.CourseGetAll)
	app.Get("/api/v1/recommended", h.CourseBasicToLearn)

	app.Post("/admin/login", func(c *fiber.Ctx) error {
		req := new(service.LoginRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}

		token, err := authService.Login(*req)
		if err != nil {
			return err
		}

		cookie := new(fiber.Cookie)
		cookie.Name = "token"
		cookie.Value = token

		c.Cookie(cookie)

		return c.Redirect("/admin")
	})

	app.Get("/admin", middlewareService.IsAuthenticatedFromCookie, func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(model.User)

		if !ok {
			return errors.New("cannot convert to user pointer")
		}

		courses, images, err := courseService.GetAll()
		if err != nil {
			return err
		}

		return c.Render("admin/index", fiber.Map{"User": user, "Courses": courses, "Images": images})
	})

	app.Get("/admin/courses/create", middlewareService.IsAuthenticatedFromCookie, func(c *fiber.Ctx) error {
		user, ok := c.Locals("user").(model.User)

		if !ok {
			return errors.New("cannot convert to user pointer")
		}

		return c.Render("admin/courses/create", fiber.Map{"User": user})
	})

	app.Post("/admin/courses", middlewareService.IsAuthenticatedFromCookie, func(c *fiber.Ctx) error {
		req := new(service.CreateCourseRequest)
		if err := c.BodyParser(req); err != nil {
			return err
		}

		image, err := c.FormFile("image")
		if err != nil {
			return err
		}

		imageAsFile, err := image.Open()
		if err != nil {
			return err
		}

		ctx := context.Background()
		resp, err := cld.Upload.Upload(ctx, imageAsFile, uploader.UploadParams{})
		if err != nil {
			return err
		}

		err = courseService.Create(*req, resp.PublicID)
		if err != nil {
			return err
		}

		return c.Redirect("/admin")
	})

	app.Get("/admin/login", func(c *fiber.Ctx) error {
		return c.Render("admin/login", fiber.Map{})
	})

	if err := app.Listen(":8080"); err != nil {
		log.Fatalln("Cannot listen: ", err.Error())
	}
}
