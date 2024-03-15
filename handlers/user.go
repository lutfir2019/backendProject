package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	guuid "github.com/google/uuid"
	"go.mod/database"
	"go.mod/handlers/structur"
	"go.mod/helper"
	"go.mod/model"
	"gorm.io/gorm"
)

func CreateUser(c *fiber.Ctx) error {
	// Ambil token dari fiber.Ctx.Locals()
	token, ok := c.Locals("token").(*jwt.Token)
	if !ok {
		return helper.ResponseBasic(c, 400, "Failed to get token from context")
	}

	// user := c.Locals("token").(structur.Token)

	fmt.Println(token)

	json := new(structur.SliceUserRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, "Invalid JSON", err)
	}

	// cek apakah toko tersedia
	db := database.DB
	shop := Shop{}
	queryShop := Shop{Spcd: json.Spcd}
	err := db.First(&shop, &queryShop).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, "Shop not found", err)
	}

	pass := helper.HashAndSalt([]byte(json.Pass))
	new := User{
		UID:       guuid.New(),
		Nam:       json.Nam,
		Unm:       json.Unm,
		Pass:      pass,
		Rlcd:      json.Rlcd,
		Rlnm:      json.Rlnm,
		Almt:      json.Almt,
		Gdr:       json.Gdr,
		Pn:        json.Pn,
		Spcd:      json.Spcd,
		Spnm:      json.Spnm,
		ShopRefer: shop.SID,
	}

	found := User{}
	query := User{Unm: new.Unm}
	err = db.First(&found, &query).Error
	if err != gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 400, "User already exists", err)
	}

	err = db.Create(&new).Error
	if err != nil {
		return helper.ResponsError(c, 500, "Invalid query database", err)
	}

	return helper.ResponsSuccess(c, 200, "Success create user", found, 1, 10, 1)
}

func GetUsers(c *fiber.Ctx) error {
	json := new(structur.SizeGetDataRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, "Invalid JSON", err)
	}

	// Set default value if not set in the request page
	if json.Page < 1 {
		json.Page = 1
	}
	if json.PageSize < 1 {
		json.PageSize = 10
	}
	offset := (json.Page - 1) * json.PageSize

	db := database.DB
	Users := []User{}
	db.Model(&model.User{}).Count(&TotalItems)
	db.Model(&model.User{}).Order("ID DESC").Offset(offset).Limit(json.PageSize).Find(&Users)

	return helper.ResponsSuccess(c, 200, "Succes get data user", Users, TotalItems, json.PageSize, json.Page)
}

func GetUserByUnm(c *fiber.Ctx) error {
	param := c.Params("unm")

	db := database.DB
	user := User{}
	query := User{Unm: param}
	err := db.First(&user, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 200, "User not found", err)
	}

	return helper.ResponsSuccess(c, 200, "Succes get data user by username", user, 1, 10, 1)
}

func UpdateUserByUnm(c *fiber.Ctx) error {
	param := c.Params("unm")

	json := new(structur.SliceUserRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, "Invalid JSON", err)
	}

	db := database.DB
	user := User{}
	query := User{Unm: param}
	err := db.First(&user, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, "User not found", err)
	}

	// cek apakah ada field yang di update
	if json.Nam != "" {
		user.Nam = json.Nam
	}
	if json.Almt != "" {
		user.Almt = json.Almt
	}
	if json.Gdr != "" {
		user.Gdr = json.Gdr
	}
	if json.Rlcd != "" {
		user.Rlcd = json.Rlcd
	}
	if json.Rlnm != "" {
		user.Rlnm = json.Rlnm
	}
	if json.Pn != "" {
		user.Pn = json.Pn
	}
	if json.Pn != "" {
		user.Pn = json.Pn
	}

	db.Save(user)
	return helper.ResponsSuccess(c, 200, "Succes get data user", User{}, 1, 10, 1)
}

func DeleteByUnm(c *fiber.Ctx) error {
	param := c.Params("unm")

	db := database.DB
	user := User{}
	query := User{Unm: param}
	err := db.First(&user, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, "User not found", err)
	}

	// DB.Model(&found).Association("Sessions").Delete()
	db.Model(&user).Association("Products").Delete()
	db.Delete(&user)

	username := fmt.Sprintf("Success delete user data %s", user.Unm)
	return helper.ResponseBasic(c, 200, username)
}

func ChangePassword(c *fiber.Ctx) error {
	json := new(structur.ChangePasswordRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, "Invalid JSON", err)
	}

	db := database.DB
	user := User{}
	query := User{Unm: json.Unm}
	err := db.First(&user, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, "User not found", err)
	}

	if !helper.ComparePasswords(user.Pass, []byte(json.Pass)) {
		return helper.ResponsError(c, 400, "Invalid Password", err)
	}

	user.Pass = helper.HashAndSalt([]byte(json.NewPass))
	db.Save(&user)
	
	return helper.ResponseBasic(c, 200, "Success change password")
}
