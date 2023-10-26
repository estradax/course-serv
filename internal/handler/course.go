package handler

import (
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CourseGetAll(c *fiber.Ctx) error {
	return c.SendString("Lick me daddy")
}
