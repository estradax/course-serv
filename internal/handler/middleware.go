package handler

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/estradax/course-serv/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func (h *Handler) Authenticated(c *fiber.Ctx) error {
	authToken := c.Get("Authorization")
	splitToken := strings.Split(authToken, " ")

	if len(splitToken) != 2 {
		log.Println("token: ", authToken)
		return errors.New("invalid authorization token")
	}

	jwtToken := splitToken[1]

	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method jwt")
		}

		return h.Secret, nil
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
		return errCantQueryFirstUser
	}

	c.Locals("user", user)

	return c.Next()
}
