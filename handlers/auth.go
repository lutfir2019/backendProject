package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"

	guuid "github.com/google/uuid"

	"go.mod/database"
	"go.mod/handlers/structur"
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

const InvalidJson = "Invalid JSON"
const NotFoundProduct = "Product not found"
const NotFoundUser = "User not found"
const NotFoundShop = "Shop not found"

func Login(c *fiber.Ctx) error {
	json := new(structur.LoginRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, InvalidJson, err)
	}

	db := database.DB
	found := User{}
	query := User{Unm: json.Unm}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, "Username not found", err)
	}
	if !helper.ComparePasswords(found.Pass, []byte(json.Pass)) {
		return helper.ResponseBasic(c, 400, "Invalid Password")
	}

	newSession := Session{
		Sessionid: guuid.New(),
		Expires:   helper.SessionExpires(),
		UserRefer: found.UID,
	}

	err = db.Create(&newSession).Error
	if err != nil {
		return helper.ResponsError(c, 500, "Invalid query database", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "__s",
		Expires:  helper.SessionExpires(),
		Value:    newSession.Sessionid.String(),
		HTTPOnly: true,
	})

	return c.Status(200).JSON(fiber.Map{
		"message": "Logged in successfully",
		"token":   newSession.Sessionid.String(),
		"data":    found,
	})
}

func Logout(c *fiber.Ctx) error {
	json := new(Session)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, InvalidJson, err)
	}

	db := database.DB
	session := Session{}
	query := Session{Sessionid: json.Sessionid}
	err := db.First(&session, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponseBasic(c, 200, "Logged out successfully")
	}

	err = db.Delete(&session).Error
	if err != nil {
		return helper.ResponsError(c, 500, "Failed to delete session", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:    "__s",
		Expires: time.Now().Add(-1 * time.Second),
	})

	return helper.ResponseBasic(c, 200, "Logged out successfully")
}
