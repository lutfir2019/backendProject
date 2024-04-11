package handlers

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"go.mod/database"
	"go.mod/handlers/structur"
	"go.mod/helper"
	"go.mod/middleware"
	"go.mod/model"
	"gorm.io/gorm"
)

func CreateUser(c *fiber.Ctx) error {
	json := new(structur.SliceUserRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, InvalidJson, err)
	}

	localUser := c.Locals("user").(map[string]interface{})
	if localUser == nil || localUser["role"].(string) == "ROLE-3" {
		return helper.ResponseBasic(c, 403, "Forbiden")
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

	user := User{}
	query := User{Unm: new.Unm, Spcd: new.Spcd}
	err = db.First(&user, &query).Error
	if err != gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 400, "Username already exists", err)
	}

	err = db.Create(&new).Error
	if err != nil {
		return helper.ResponsError(c, 500, "Invalid query database", err)
	}

	return helper.ResponsSuccess(c, 200, "Success create user", new, 1, 10, 1)
}

func GetUsers(c *fiber.Ctx) error {
	json := new(structur.SizeGetDataRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, InvalidJson, err)
	}

	localUser := c.Locals("user").(map[string]interface{})
	if localUser == nil || localUser["role"].(string) == "ROLE-3" {
		return helper.ResponseBasic(c, 403, "Forbiden")
	}
	spcd, ok := localUser["shopcode"].(string)
	if !ok || spcd == "" {
		return helper.ResponseBasic(c, 404, "Invalid Shop code")
	}
	role, ok := localUser["role"].(string)
	if !ok || role == "" {
		return helper.ResponseBasic(c, 404, "Invalid Role code")
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

	query := db.Model(&model.User{}).Order("ID DESC").Where("LOWER(nam) LIKE ? AND LOWER(spcd) LIKE ?", "%"+strings.ToLower(json.Nam)+"%", "%"+strings.ToLower(json.Spcd)+"%")

	if role != "ROLE-1" {
		query = query.Where(&model.User{Spcd: spcd})
	}

	// Eksekusi query dan simpan hasilnya di dalam Users
	query.Count(&TotalItems).Offset(offset).Limit(json.PageSize).Find(&Users)

	return helper.ResponsSuccess(c, 200, "Succes get data user", Users, TotalItems, json.PageSize, json.Page)
}

func GetUserByUnm(c *fiber.Ctx) error {
	param := c.Params("unm")

	localUser := c.Locals("user").(map[string]interface{})
	spcd, ok := localUser["shopcode"].(string)
	if !ok || spcd == "" {
		return helper.ResponseBasic(c, 404, "Invalid Shop code")
	}
	role, ok := localUser["role"].(string)
	if !ok || role == "" {
		return helper.ResponseBasic(c, 404, "Invalid Role code")
	}

	db := database.DB
	user := User{}
	query := User{Unm: param}
	if role != "ROLE-1" {
		query = User{Unm: param, Spcd: spcd}
	}
	err := db.First(&user, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 200, NotFoundUser, err)
	}

	return helper.ResponsSuccess(c, 200, "Succes get data user by username", user, 1, 10, 1)
}

func UpdateUserByUnm(c *fiber.Ctx) error {
	json := new(structur.SliceUserRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, InvalidJson, err)
	}

	db := database.DB
	user := User{}
	query := User{Unm: json.Unm, Spcd: json.Spcd}
	err := db.First(&user, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, NotFoundUser, err)
	}

	err = db.Model(&model.User{}).Where("unm =? AND spcd =?", json.Unm, json.Spcd).Updates(json).Error
	if err != nil {
		return helper.ResponsError(c, 500, "Invalid query databsae", err)
	}
	return helper.ResponsSuccess(c, 200, "Succes update data user", User{}, 1, 10, 1)
}

func DeleteByUnm(c *fiber.Ctx) error {
	json := new(structur.SliceUserRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, InvalidJson, err)
	}

	if err := middleware.DenyForStaff(c); err != nil {
		return err // Mengembalikan respons error dari middleware
	}

	db := database.DB
	user := User{}
	query := User{Unm: json.Unm, Spcd: json.Spcd}
	err := db.First(&user, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, NotFoundUser, err)
	}

	db.Model(&user).Association("Sessions").Delete()
	db.Model(&user).Association("Products").Delete()
	db.Delete(&user)

	username := fmt.Sprintf("Success delete user data %s", user.Unm)
	return helper.ResponseBasic(c, 200, username)
}

func ChangePassword(c *fiber.Ctx) error {
	json := new(structur.ChangePasswordRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, InvalidJson, err)
	}

	db := database.DB
	user := User{}
	query := User{Unm: json.Unm, Spcd: json.Spcd}
	err := db.First(&user, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, NotFoundUser, err)
	}

	if !helper.ComparePasswords(user.Pass, []byte(json.Pass)) {
		return helper.ResponsError(c, 400, "Invalid Password", err)
	}

	user.Pass = helper.HashAndSalt([]byte(json.NewPass))
	db.Save(&user)

	return helper.ResponseBasic(c, 200, "Success change password")
}
