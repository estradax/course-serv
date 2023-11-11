package middleware

import (
	"errors"
	"strconv"

	"github.com/estradax/course-serv/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type Middleware struct {
	DB     *gorm.DB
	Secret []byte
}

func New(db *gorm.DB, secret []byte) *Middleware {
	return &Middleware{
		DB: db, 
		Secret: secret,
	}
}

func JwtToSub(jwtToken string, secret []byte) (uint64, error) {
	token, err := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method jwt")
		}

		return secret, nil
	})
	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("token invalid")
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseUint(sub, 10, 32)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *Middleware) IsAuthenticatedFromCookie(c *fiber.Ctx) error {
	jwtToken := c.Cookies("token")

	id, err := JwtToSub(jwtToken, m.Secret)
	if err != nil {
		return c.Redirect("/admin/login")
	}

	user := model.User{
		Model: gorm.Model{
			ID: uint(id),
		},
	}

	result := m.DB.First(&user)
	if result.Error != nil {
		return c.Redirect("/admin/login")
	}

	c.Locals("user", user)

	return c.Next()

}
