package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"go.mod/database"
	"go.mod/handlers/structur"
	"go.mod/helper"
)

func CreateTransaction(c *fiber.Ctx) error {
	json := new(structur.CreateTransactionRequest)
	if err := c.BodyParser(json); err != nil {
		return helper.ResponsError(c, fiber.StatusBadRequest, InvalidJson, err)
	}

	db := database.DB
	tx := db.Begin() // Memulai transaksi

	for _, item := range json.Data {
		product := Product{}
		query := Product{Pcd: item.ProductCode, Spcd: item.ShopCode}
		err := tx.First(&product, &query).Error
		if err == gorm.ErrRecordNotFound {
			tx.Rollback() // Mengembalikan transaksi karena terjadi kesalahan
			return helper.ResponseBasic(c, 400, fmt.Sprintf("The product with Code %s is not found.", item.ProductCode))
		}

		// Mencari toko berdasarkan kode yang diberikan dalam JSON
		shop := Shop{}
		queryShop := Shop{Spcd: item.ShopCode}
		err = tx.First(&shop, &queryShop).Error
		if err != nil {
			tx.Rollback() // Mengembalikan transaksi karena terjadi kesalahan
			return helper.ResponseBasic(c, 400, fmt.Sprintf("Invalid code from %s.", item.ShopCode))
		}

		newTransaction := Transaction{
			Total:        item.Total,
			Price:        item.Price,
			Quantity:     item.Quantity,
			Catnm:        item.Catnm,
			Catcd:        item.Catcd,
			ProductRefer: product.PID,
			ShopRefer:    shop.SID,
		}

		// Create the product
		err = tx.Create(&newTransaction).Error
		if err != nil {
			tx.Rollback() // Mengembalikan transaksi karena terjadi kesalahan
			return helper.ResponseBasic(c, 500, "Invalid query database")
		}

		// Update product quantity
		product.Qty -= item.Quantity
		err = tx.Save(&product).Error
		if err != nil {
			tx.Rollback() // Rollback transaction due to error
			return helper.ResponseBasic(c, 500, "Failed to update product quantity")
		}
	}

	tx.Commit() // Melakukan commit transaksi jika tidak ada kesalahan

	return helper.ResponseBasic(c, 200, "Transaction successfully")
}
