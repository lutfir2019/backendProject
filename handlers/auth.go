package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	guuid "github.com/google/uuid"
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
	db         = database.DB
	SecretKey  = []byte("IUAacnfkjdxMJXO;ALSKZXCSOIGAJFMDSKAMsijkd[0ANUG0[")
	totalItems int64
)

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Unm  string `json:"unm"`
		Pass string `json:"pass"`
	}

	db := database.DB
	json := new(LoginRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

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

	isAdmin := found.Rlcd != "ROLE-3"

	// Create the Claims
	claims := jwt.MapClaims{
		"name":  found.Unm,
		"admin": isAdmin,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString(SecretKey)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Failed to sign session",
			"details": err,
		})
	}

	session := Session{
		UserRefer: found.UID,
		Expires:   helper.SessionExpires(),
		Sessionid: guuid.New(),
		Token:     "Bearer " + t,
	}

	db.Create(&session)
	// c.Cookie(&fiber.Cookie{
	// 	Name:     "sessionid",
	// 	Expires:  helper.SessionExpires(),
	// 	Value:    session.Token,
	// 	HTTPOnly: true,
	// })

	return c.Status(200).JSON(fiber.Map{
		"message": "Sucessfully",
		"data":    session,
	})
}

func Logout(c *fiber.Ctx) error {
	db := database.DB
	json := new(Session)
	if err := c.BodyParser(json); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}
	session := Session{}
	query := Session{Sessionid: json.Sessionid}
	err := db.First(&session, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"message": "Session not found",
		})
	}
	db.Delete(&session)
	c.ClearCookie("sessionid")
	return c.Status(200).JSON(fiber.Map{
		"message": "sucess",
	})
}
