package service

import (
	"github.com/estradax/course-serv/internal/model"
	"gorm.io/gorm"
)

type Course struct {
	DB     *gorm.DB
	Secret []byte
}

func NewCourseService(db *gorm.DB, secret []byte) *Course {
	return &Course{
		DB:     db,
		Secret: secret,
	}
}

func (s *Course) GetAll() ([]model.Course, error) {
	courses := new([]model.Course)

	result := s.DB.Find(courses)
	if result.Error != nil {
		return []model.Course{}, result.Error
	}

	return *courses, nil
}
