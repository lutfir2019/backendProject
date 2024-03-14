package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"go.mod/database"
	"go.mod/handlers/structur"
	"go.mod/helper"
	"go.mod/model"
	"gorm.io/gorm"
)

func CreateShop(c *fiber.Ctx) error {
	db := database.DB
	json := new(structur.SliceShopRequest)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	newShop := Shop{
		Spnm: json.Spnm,
		Almt: json.Almt,
		SID:  guuid.New(), // generate UUID for shop code
		Spcd: helper.GenerateCode(json.Spnm),
	}

	found := Shop{}
	query := Shop{Spnm: newShop.Spnm}
	err := db.First(&found, &query).Error
	if err != gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Shop already exists",
		})
	}

	err = db.Create(&newShop).Error
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Success create shop data",
	})
}

func GetShops(c *fiber.Ctx) error {
	db := database.DB

	// Kebutuhan pagination
	page, err := strconv.Atoi(c.Query("page", "1")) // Ambil nomor halaman dari query parameter, defaultnya 1
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("page_size", "10")) // Ambil jumlah item per halaman dari query parameter, defaultnya 10
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit // Hitung offset berdasarkan nomor halaman

	var totalItems int64
	Shops := []Shop{}

	// Hitung total item (untuk kebutuhan pagination)
	db.Model(&model.Shop{}).Count(&totalItems)

	// Ambil data toko dengan pagination
	db.Model(&model.Shop{}).Order("s_id").Offset(offset).Limit(limit).Find(&Shops)

	return c.Status(200).JSON(fiber.Map{
		"message": "Success get shop data",
		"data": fiber.Map{
			"data":       Shops,
			"pagination": helper.SetPagination(totalItems, limit, page),
		},
	})
}

func GetShopByCode(c *fiber.Ctx) error {
	db := database.DB
	param := c.Params("sid")
	sid, err := guuid.Parse(param)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request params sid",
		})
	}

	// Kebutuhan pagination
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(c.Query("page_size", "10"))
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	var totalItems int64

	shop := Shop{}
	query := Shop{SID: sid}
	err = db.First(&shop, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "Shop not found",
			"data": fiber.Map{
				"data":       []Shop{},
				"pagination": helper.SetPagination(totalItems, limit, page),
			},
		})
	}

	db.Model(&model.Shop{}).Where(&query).Count(&totalItems)
	db.Model(&model.Shop{}).Where(&query).Order("SID asc").Offset(offset).Limit(limit).Find(&shop)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Success get data",
		"data": fiber.Map{
			"data":       shop,
			"pagination": helper.SetPagination(totalItems, limit, page),
		},
	})
}

func UpdateShop(c *fiber.Ctx) error {
	db := database.DB
	param := c.Params("sid")
	sid, err := guuid.Parse(param)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid request parameter",
		})
	}

	json := new(structur.SliceShopRequest)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	found := Shop{}
	query := Shop{
		SID: sid,
	}

	err = db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "Shop not found",
		})
	}

	if json.Spnm != "" {
		found.Spnm = json.Spnm
	}
	if json.Almt != "" {
		found.Almt = json.Almt
	}

	db.Save(&found)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Sucsess update shop data",
	})
}

func DeleteShop(c *fiber.Ctx) error {
	db := database.DB
	param := c.Params("sid")
	sid, err := guuid.Parse(param)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid SPCD parameter",
		})
	}

	found := Shop{}
	query := Shop{
		SID: sid,
	}

	err = db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Shop not found",
		})
	}

	db.Delete(&found)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Sucess  delete shop data",
	})
}
