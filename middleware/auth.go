package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mod/helper"
)

// Middleware untuk memeriksa token JWT
func Authenticated(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("__s") != "" {
		tokenString = c.Cookies("__s")
	}
	if tokenString == "" {
		return helper.ResponseBasic(c, fiber.StatusUnauthorized, "You are not logged in")
	}

	return helper.ParseSessionId(c, tokenString)
}

func DenyForStaff(c *fiber.Ctx) error {
	localUser := c.Locals("user").(map[string]interface{})
	if localUser == nil || localUser["role"].(string) == "ROLE-3" {
		return helper.ResponseBasic(c, 403, "Forbiden")
	}
	fmt.Println(localUser)
	return c.Next()
}
