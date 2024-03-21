package handlers

import (
	"github.com/gofiber/fiber/v2"
	"go.mod/database"
	"go.mod/helper"
)

// Struct untuk menyimpan hasil query
type Result struct {
	Category string
	Total    int
}

func GetTransactionByPeriod(c *fiber.Ctx) error {

	db := database.DB
	// transaksiPenjualan := db.Model(&[]model.Transaction{}).Order("ID DESC")

	// // Ambil data transaksi penjualan untuk hari ini
	// db.Where("DATE(transaction_date) = ?", time.Now().Format("2006-01-02")).Find(&transaksiPenjualan)

	// // Ambil data transaksi penjualan untuk minggu ini
	// db.Where("WEEK(transaction_date) = ? AND YEAR(transaction_date) = ?", time.Now().Weekday(), time.Now().Year()).Find(&transaksiPenjualan)

	// // Ambil data transaksi penjualan untuk bulan ini
	// db.Where("MONTH(transaction_date) = ? AND YEAR(transaction_date) = ?", time.Now().Month(), time.Now().Year()).Find(&transaksiPenjualan)

	// // Ambil data transaksi penjualan untuk tahun ini
	// db.Where("YEAR(transaction_date) = ?", time.Now().Year()).Find(&transaksiPenjualan)

	// return helper.ResponseBasic(c, 200, "Hallo world!")
	var results []Result

	// Query untuk mengelompokkan transaksi berdasarkan kategori produk dan menghitung jumlah penjualan untuk setiap kategori
	if err := db.Table("transactions").
		Select("products.catnm AS category, SUM(transactions.quantity) AS total").
		Joins("JOIN products ON transactions.product_refer = products.p_id").
		Group("products.catnm").
		Order("total DESC").
		Find(&results).Error; err != nil {
		return helper.ResponsError(c, 500, "Invalid query", err)
	}

	if len(results) == 0 {
		return helper.ResponseBasic(c, 400, "Not found")
	}

	// Kembalikan kategori dengan penjualan terbanyak
	bestCategory := results[0].Category
	bestTotal := results[0].Total

	return c.Status(200).JSON(fiber.Map{
		"bestCategory": bestCategory,
		"bestTotal":    bestTotal,
	})
}
