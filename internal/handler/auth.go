package handler

import (
	"errors"
	"fmt"

	"github.com/estradax/course-serv/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var errCantQueryFirstUser = errors.New("can't query first user")

type Handler struct {
	DB     *gorm.DB
	Secret []byte
}

func NewHandler(db *gorm.DB, secret []byte) *Handler {
	return &Handler{
		DB:     db,
		Secret: secret,
	}
}

func (h *Handler) Profile(c *fiber.Ctx) error {
	user, ok := c.Locals("user").(model.User)

	if !ok {
		return errors.New("cannot convert to user pointer")
	}

	return c.JSON(fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

type registerRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(c *fiber.Ctx) error {
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

	token, err := createToken(jwt.RegisteredClaims{
		Subject: fmt.Sprintf("%v", user.ID),
	}, h.Secret)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(c *fiber.Ctx) error {
	p := new(loginRequest)
	if err := c.BodyParser(p); err != nil {
		return err
	}

	user := model.User{}

	result := h.DB.Where("email = ?", p.Email).First(&user)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected != 1 {
		return errCantQueryFirstUser
	}

	token, err := createToken(jwt.RegisteredClaims{
		Subject: fmt.Sprintf("%v", user.ID),
	}, h.Secret)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{
		"token": token,
	})
}

func createToken(claims jwt.RegisteredClaims, secret []byte) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := t.SignedString(secret)
	if err != nil {
		return "", err
	}

	return token, nil
}
