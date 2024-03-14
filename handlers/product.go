package handlers

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mod/database"
	"go.mod/handlers/structur"
	"go.mod/helper"
	"go.mod/model"
	"gorm.io/gorm"
)

func CreateProduct(c *fiber.Ctx) error {
	json := new(structur.CreateProductRequest)
	if err := c.BodyParser(json); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "Invalid JSON",
		})
	}

	db := database.DB
	// Loop through the products in the request
	count := 0
	for _, item := range json.Data {
		count++
		// Check if the product already exists
		found := Product{}
		query := Product{Pcd: item.Pcd}
		err := db.First(&found, &query).Error
		if err != gorm.ErrRecordNotFound {
			return c.Status(400).JSON(fiber.Map{
				"message": fmt.Sprintf("The product with Code %s is already registered.", item.Pcd),
			})
		}
		if count == len(json.Data) {
			for _, item := range json.Data {
				count++
				// Mencari toko berdasarkan kode yang diberikan dalam JSON
				shop := Shop{}
				queryShop := Shop{Spcd: item.Spcd}
				err := db.First(&shop, &queryShop).Error
				if err != nil {
					return c.Status(400).JSON(fiber.Map{
						"message": fmt.Sprintf("Invalid code from %s.", item.Spnm),
					})
				}
				if count == len(json.Data) {
					for _, item := range json.Data {
						shop := Shop{}
						queryShop := Shop{Spcd: item.Spcd}
						db.First(&shop, &queryShop)
						newProduct := Product{
							Pnm:       item.Pnm,
							Pcd:       item.Pcd,
							Qty:       item.Qty,
							Price:     item.Price,
							Catcd:     item.Catcd,
							Catnm:     item.Catnm,
							Spcd:      item.Spcd,
							Spnm:      item.Spnm,
							ShopRefer: shop.SID,
						}
						// Create the product
						err = db.Create(&newProduct).Error
						if err != nil {
							return c.Status(500).JSON(fiber.Map{
								"message": "Invalid query database",
							})
						}
					}
				}
			}
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Success create product data",
	})
}

func GetProducts(c *fiber.Ctx) error {
	db := database.DB

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
	Products := []Product{}

	db.Model(&model.Product{}).Count(&totalItems)
	db.Model(&model.Product{}).Order("p_id DESC").Offset(offset).Limit(limit).Find(&Products)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Success get product data",
		"data": fiber.Map{
			"data":       Products,
			"pagination": helper.SetPagination(totalItems, limit, page),
		},
	})
}
func GetProductByCode(c *fiber.Ctx) error {
	db := database.DB
	pcd := c.Params("pcd")

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

	product := Product{}
	query := Product{Pcd: pcd}
	err = db.First(&product, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "Product not found",
			"data": c.JSON(fiber.Map{
				"data":       []Product{},
				"pagination": helper.SetPagination(totalItems, limit, page),
			}),
		})
	}

	db.Model(&model.Product{}).Where(&query).Count(&totalItems)
	db.Model(&model.Product{}).Where(&query).Order("PID asc").Offset(offset).Limit(limit)
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "Success get product by code",
		"data": fiber.Map{
			"data": product,
			"pagination": fiber.Map{
				"totalItems":  totalItems,
				"totalPages":  int(totalItems) / limit,
				"currentPage": page,
				"perPage":     limit,
			},
		},
	})
}

func UpdateProduct(c *fiber.Ctx) error {

	db := database.DB
	user := c.Locals("user").(User)
	json := new(structur.SliceProductRequest)
	if err := c.BodyParser(json); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Invalid JSON",
		})
	}

	pcd := c.Params("pcd")
	found := Product{}
	query := Product{
		Pcd:       pcd,
		ShopRefer: user.ShopRefer,
	}

	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"message": "Product not found",
		})
	}
	if json.Pnm != "" {
		found.Pnm = json.Pnm
	}
	if json.Qty != 0 {
		found.Qty = json.Qty
	}
	if json.Price != 0 {
		found.Price = json.Price
	}
	if json.Catnm != "" {
		found.Catnm = json.Catnm
	}
	if json.Catcd != "" {
		found.Catcd = json.Catcd
	}
	db.Save(&found)
	return c.Status(200).JSON(fiber.Map{
		"message": "Success update product data",
	})
}
func DeleteProduct(c *fiber.Ctx) error {
	db := database.DB
	user := c.Locals("user").(User)
	pcd := c.Params("pcd")
	found := Product{}
	query := Product{
		Pcd:       pcd,
		ShopRefer: user.ShopRefer,
	}

	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(404).JSON(fiber.Map{
			"message": "Product not found",
		})
	}
	db.Delete(&found)
	return c.Status(200).JSON(fiber.Map{
		"message": "Success delete product data",
	})
}
