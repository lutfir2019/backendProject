package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	guuid "github.com/google/uuid"
	"go.mod/database"
	"go.mod/handlers/structur"
	"go.mod/helper"
	"go.mod/middleware"
	"go.mod/model"
	"gorm.io/gorm"
)

func CreateProduct(c *fiber.Ctx) error {
    json := new(structur.CreateProductRequest)
    if err := c.BodyParser(json); err != nil {
        return helper.ResponseBasic(c, 400, "Invalid JSON")
    }

    if err := middleware.DenyForStaff(c); err != nil {
        return err // Mengembalikan respons error dari middleware
    }

    db := database.DB
    tx := db.Begin() // Memulai transaksi

    existingProducts := make(map[string]bool)

    for _, item := range json.Data {
        // Check if the product already exists
        found := Product{}
        query := Product{Pcd: item.Pcd}
        err := tx.First(&found, &query).Error
        if err != gorm.ErrRecordNotFound {
            tx.Rollback() // Mengembalikan transaksi karena terjadi kesalahan
            return helper.ResponseBasic(c, 400, fmt.Sprintf("The product with Code %s is already registered.", item.Pcd))
        }

        // Check if the product code already exists in the current batch
        if existingProducts[item.Pcd] {
            tx.Rollback() // Mengembalikan transaksi karena terjadi kesalahan
            return helper.ResponseBasic(c, 400, fmt.Sprintf("Duplicate product code found: %s", item.Pcd))
        }

        existingProducts[item.Pcd] = true

        // Mencari toko berdasarkan kode yang diberikan dalam JSON
        shop := Shop{}
        queryShop := Shop{Spcd: item.Spcd}
        err = tx.First(&shop, &queryShop).Error
        if err != nil {
            tx.Rollback() // Mengembalikan transaksi karena terjadi kesalahan
            return helper.ResponseBasic(c, 400, fmt.Sprintf("Invalid code from %s.", item.Spnm))
        }

        newProduct := Product{
			PID:       guuid.New(),
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
        err = tx.Create(&newProduct).Error
        if err != nil {
            tx.Rollback() // Mengembalikan transaksi karena terjadi kesalahan
            return helper.ResponseBasic(c, 500, "Invalid query database")
        }
    }

    tx.Commit() // Melakukan commit transaksi jika tidak ada kesalahan

    return helper.ResponseBasic(c, 200, "Success create product data")
}


func GetProducts(c *fiber.Ctx) error {
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
	Products := []Product{}

	db.Model(&model.Product{}).Count(&TotalItems)
	db.Model(&model.Product{}).Order("ID DESC").Offset(offset).Limit(json.PageSize).Find(&Products)

	return helper.ResponsSuccess(c, 200, "Success get product data", Products, TotalItems, json.PageSize, json.Page)
}
func GetProductByCode(c *fiber.Ctx) error {
	param := c.Params("pcd")

	db := database.DB
	product := Product{}
	query := Product{Pcd: param}
	err := db.First(&product, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, "Product not found", err)
	}

	return helper.ResponsSuccess(c, 200, "Success get product by code", product, 1, 10, 1)
}

func UpdateProductByCode(c *fiber.Ctx) error {
	param := c.Params("pcd")

	if err := middleware.DenyForStaff(c); err != nil {
		return err // Mengembalikan respons error dari middleware
	}

	json := new(structur.SliceProductRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, 400, "Invalid JSON", err)
	}

	user := c.Locals("user")

	fmt.Println(user)

	db := database.DB
	product := Product{}
	query := Product{Pcd: param}
	err := db.First(&product, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, "Product not found", err)
	}

	if json.Pnm != "" {
		product.Pnm = json.Pnm
	}
	if json.Qty != 0 {
		product.Qty = json.Qty
	}
	if json.Price != 0 {
		product.Price = json.Price
	}
	if json.Catnm != "" {
		product.Catnm = json.Catnm
	}
	if json.Catcd != "" {
		product.Catcd = json.Catcd
	}

	db.Save(&product)

	return helper.ResponsSuccess(c, 200, "Success update product data", Product{}, 1, 10, 1)
}

func DeleteProduct(c *fiber.Ctx) error {
	param := c.Params("pcd")

	if err := middleware.DenyForStaff(c); err != nil {
		return err // Mengembalikan respons error dari middleware
	}

	db := database.DB
	found := Product{}
	query := Product{Pcd: param}

	err := db.First(&found, &query).Error
	if err == gorm.ErrRecordNotFound {
		return helper.ResponsError(c, 404, "Product not found", err)
	}

	db.Delete(&found)

	return helper.ResponseBasic(c, 200, fmt.Sprintf("Success delete product with code %s", param))
}
