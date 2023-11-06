package model

import "gorm.io/gorm"

type Course struct {
	gorm.Model
	Title         string
	Description   string
	Price         int32
	ImagePublicID string
}
