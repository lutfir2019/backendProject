package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mod/handlers"
	"go.mod/helper"
)

// Middleware untuk memeriksa token JWT
func Authenticated(c *fiber.Ctx) error {
	// Ambil token dari header Authorization
	authHeader := c.Get("Authorization")

	// Periksa jika Authorization header tidak ada
	if authHeader == "" {
		return helper.ResponseBasic(c, 401, "Authorization header missing")
	}

	// Pecah string menjadi bagian "Bearer" dan token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return helper.ResponseBasic(c, 401, "Invalid Authorization header format")
	}

	tokenString := parts[1]

	// Validasi token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi metode penandatanganan token
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return handlers.SecretKey, nil
	})

	if err != nil {
		return helper.ResponseBasic(c, fiber.StatusUnauthorized, "Unauthorized")
	}

	// Periksa apakah token valid
	if !token.Valid {
		return helper.ResponseBasic(c, fiber.StatusUnauthorized, "Invalid token")
	}

	return helper.ParseJwtToken(c, token)
}
