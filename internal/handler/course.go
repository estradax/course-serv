package handler

import (
	"github.com/estradax/course-serv/internal/model"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CourseGetAll(c *fiber.Ctx) error {
	users := new([]model.Course)

	result := h.DB.Find(users)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(fiber.Map{
		"Courses": users,
	})
}
