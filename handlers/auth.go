package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"

	// guuid "github.com/google/uuid"

	"go.mod/database"
	"go.mod/helper"
	"go.mod/model"
	"gorm.io/gorm"
)

type User model.User
type Session model.Session
type Product model.Product
type Shop model.Shop

var (
	SecretKey  = []byte("IUAacnfkjdxMJXO;ALSKZXCSOIGAJFMDSKAMsijkd[0ANUG0[")
	TotalItems int64
)

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Unm  string `json:"unm"`
		Pass string `json:"pass"`
	}

	json := new(LoginRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	db := database.DB
	found := User{}
	query := User{Unm: json.Unm}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"message": "Username not found",
		})
	}
	if !helper.ComparePasswords(found.Pass, []byte(json.Pass)) {
		return c.Status(401).JSON(fiber.Map{
			"message": "Invalid Password",
		})
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"username": found.Unm,
		"role":     found.Rlcd,
		"exp":      helper.SessionExpires().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return helper.ResponsError(c, 400, "Failed to create token", err)
	}

	// session := Session{
	// 	UserRefer: found.UID,
	// 	Expires:   helper.SessionExpires(),
	// 	Sessionid: guuid.New(),
	// 	Token:     "Bearer " + t,
	// }

	// db.Create(&session)
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Expires:  helper.SessionExpires(),
		Value:    t,
		HTTPOnly: true,
	})

	return c.Status(200).JSON(fiber.Map{
		"message": "Sucessfully",
		"token":   t,
	})
}

func Logout(c *fiber.Ctx) error {
	c.ClearCookie("token")
	return c.Status(200).JSON(fiber.Map{
		"message": "Logout sucessfully",
	})
}
