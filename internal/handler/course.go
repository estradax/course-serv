package handler

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/estradax/course-serv/internal/model"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) CourseGetAll(c *fiber.Ctx) error {
	courses := new([]model.Course)

	result := h.DB.Find(courses)
	if result.Error != nil {
		return result.Error
	}

	return c.JSON(fiber.Map{
		"Courses": courses,
	})
}

func (h *Handler) CourseBasicToLearn(c *fiber.Ctx) error {
	courses := new([]model.Course)

	result := h.DB.Limit(3).Order("created_at ASC").Find(courses)
	if result.Error != nil {
		return result.Error
	}

	images := []fiber.Map{}

	for _, course := range *courses {
		ctx := context.Background()
		resp, err := h.Cloudinary.Admin.Asset(ctx, admin.AssetParams{
			PublicID: course.ImagePublicID,
		})
		if err != nil {
			return err
		}

		images = append(images, fiber.Map{
			"ID":       course.ID,
			"ImageURL": resp.SecureURL,
		})
	}

	return c.JSON(fiber.Map{
		"courses": courses,
		"images":  images,
	})
}
