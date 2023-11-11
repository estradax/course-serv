package service

import (
	"fmt"

	"github.com/estradax/course-serv/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Auth struct {
	DB     *gorm.DB
	Secret []byte
}

func NewAuthService(db *gorm.DB, secret []byte) *Auth {
	return &Auth{
		DB:     db,
		Secret: secret,
	}
}

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (s *Auth) Login(req LoginRequest) (string, error) {
	user := model.User{}

	result := s.DB.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		return "", result.Error
	}

	token, err := createToken(jwt.RegisteredClaims{
		Subject: fmt.Sprintf("%v", user.ID),
	}, s.Secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func createToken(claims jwt.RegisteredClaims, secret []byte) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString(secret)
	if err != nil {
		return "", err
	}

	return token, nil
}
