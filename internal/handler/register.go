package handler

import (
	"errors"

	"github.com/estradax/course-serv/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterHandler struct {
	DB        *gorm.DB
	JWTSecret []byte
}

func NewRegisterHandler(db *gorm.DB, jwtSecret []byte) *RegisterHandler {
	return &RegisterHandler{
		DB:        db,
		JWTSecret: jwtSecret,
	}
}

func (h *RegisterHandler) Register(c *fiber.Ctx) error {
	p := new(registerRequest)
	if err := c.BodyParser(p); err != nil {
		return err
	}

	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(p.Password), 14)
	if err != nil {
		return err
	}

	user := model.User{
		Name:     p.Name,
		Email:    p.Email,
		Password: string(passwordBytes),
	}

	result := h.DB.Create(&user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return errors.New("error creating user")
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
	})
	token, err := t.SignedString(h.JWTSecret)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}
