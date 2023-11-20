package service

import (
	"context"
	"strconv"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/estradax/course-serv/internal/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Course struct {
	DB     *gorm.DB
	Secret []byte
	Cloudinary *cloudinary.Cloudinary
}

type CreateCourseRequest struct {
	Image string `json:"image" form:"image"`
	Title string `json:"title" form:"title"`
	Description string `json:"description" form:"description"`
	Price string `json:"price" form:"price"`
}

func NewCourseService(db *gorm.DB, secret []byte, cld *cloudinary.Cloudinary) *Course {
	return &Course{
		DB:     db,
		Secret: secret,
		Cloudinary: cld,
	}
}

func GetImages(courses []model.Course, cld *cloudinary.Cloudinary) ([]fiber.Map, error) {
	images := []fiber.Map{}

	for _, course := range courses {
		ctx := context.Background()
		resp, err := cld.Admin.Asset(ctx, admin.AssetParams{
			PublicID: course.ImagePublicID,
		})
		if err != nil {
			return nil, err
		}

		images = append(images, fiber.Map{
			"ID":       course.ID,
			"ImageURL": resp.SecureURL,
		})
	}

	return images, nil
}

func (s *Course) GetAll() ([]model.Course, []fiber.Map, error) {
	courses := new([]model.Course)

	result := s.DB.Find(courses)
	if result.Error != nil {
		return []model.Course{}, nil, result.Error
	}

	images, err := GetImages(*courses, s.Cloudinary)
	if err != nil {
		return nil, nil, err
	}

	return *courses, images, nil
}

func (s *Course) Create(req CreateCourseRequest, imagePublicID string) error {
	price, err := strconv.Atoi(req.Price)
	if err != nil {
		return err
	}

	course := model.Course{
		Title: req.Title,
		Description: req.Description,
		Price: int32(price),
		ImagePublicID: imagePublicID,
	}

	result := s.DB.Create(&course)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
