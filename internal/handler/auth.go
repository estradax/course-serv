package handler

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

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

type AuthHandler struct {
	DB        *gorm.DB
	JWTSecret []byte
}

func NewAuthHandler(db *gorm.DB, jwtSecret []byte) *AuthHandler {
	return &AuthHandler{
		DB:        db,
		JWTSecret: jwtSecret,
	}
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	authToken := c.Get("Authorization")
	splitToken := strings.Split(authToken, " ")

	if len(splitToken) != 2 {
		return errors.New("invalid authorization token")
	}

	jwtToken := splitToken[1]

	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method jwt")
		}

		return h.JWTSecret, nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("token invalid")
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(sub, 10, 32)
	if err != nil {
		return err
	}

	user := model.User{
		Model: gorm.Model{
			ID: uint(id),
		},
	}

	result := h.DB.First(&user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return errors.New("can't query first users")
	}

	return c.JSON(fiber.Map{
		"id": user.ID,
		"name": user.Name,
		"email": user.Email,
	})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
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

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Subject: fmt.Sprintf("%v", user.ID),
	})

	token, err := t.SignedString(h.JWTSecret)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}