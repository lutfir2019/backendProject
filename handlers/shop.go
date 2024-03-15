package handlers

import (
	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"go.mod/database"
	"go.mod/handlers/structur"
	"go.mod/helper"
	"go.mod/model"
	"gorm.io/gorm"
)

func CreateShop(c *fiber.Ctx) error {
	json := new(structur.SliceShopRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, "Invalid JSON", err)
	}

	newShop := Shop{
		SID:  guuid.New(), // generate UUID for shop code
		Spnm: json.Spnm,
		Spcd: helper.GenerateCode(json.Spnm),
		Almt: json.Almt,
	}

	db := database.DB
	found := Shop{}
	query := Shop{Spnm: newShop.Spnm}
	err := db.First(&found, &query).Error
	if err != gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 400, "Name shop already exists", err)
	}

	err = db.Create(&newShop).Error
	if err != nil {
		return helper.ResponsError(c, 500, "Invalid query database", err)
	}

	return helper.ResponseBasic(c, 200, "Success create shop data")

}

func GetShops(c *fiber.Ctx) error {
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
	Shops := []Shop{}

	// Hitung total item (untuk kebutuhan pagination)
	db.Model(&model.Shop{}).Count(&TotalItems)
	// Ambil data toko dengan pagination
	db.Model(&model.Shop{}).Order("ID DESC").Order("s_id").Offset(offset).Limit(json.PageSize).Find(&Shops)

	return helper.ResponsSuccess(c, 200, "Succes get data user", Shops, TotalItems, json.PageSize, json.Page)
}

func GetShopByCode(c *fiber.Ctx) error {
	param := c.Params("scd")

	json := new(structur.SizeGetDataRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, "Invalid JSON", err)
	}

	if json.Page < 1 {
		json.Page = 1
	}
	if json.PageSize < 1 {
		json.PageSize = 10
	}
	offset := (json.Page - 1) * json.PageSize

	db := database.DB
	shop := Shop{}
	query := Shop{Spcd: param}
	err := db.First(&shop, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, "Shop not found", err)
	}

	db.Model(&model.Shop{}).Where(&query).Count(&TotalItems)
	db.Model(&model.Shop{}).Where(&query).Order("ID DESC").Offset(offset).Limit(json.PageSize).Find(&shop)
	return helper.ResponsSuccess(c, 200, "Succes get shop data by code", shop, TotalItems, json.PageSize, json.Page)
}

func UpdateShop(c *fiber.Ctx) error {
	param := c.Params("scd")
	json := new(structur.SliceShopRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, "Invalid JSON", err)
	}

	db := database.DB
	found := Shop{}
	query := Shop{Spcd: param}
	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, "Shop not found", err)
	}

	if json.Spnm != "" {
		found.Spnm = json.Spnm
	}
	if json.Almt != "" {
		found.Almt = json.Almt
	}

	db.Save(&found)
	return helper.ResponseBasic(c, 200, "Sucsess update shop data")
}

func DeleteShop(c *fiber.Ctx) error {
	param := c.Params("spcd")

	db := database.DB
	found := Shop{}
	query := Shop{Spcd: param}

	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, "Shop not found", err)
	}

	db.Delete(&found)
	return helper.ResponseBasic(c, 200, "Sucess  delete shop data")
}
