package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mod/database"
	"go.mod/helper"
	"go.mod/model"
)

func GetTransactionByPeriod(c *fiber.Ctx) error {

	db := database.DB
	transaksiPenjualan := db.Model(&[]model.Transaction{}).Order("ID DESC")

	// Ambil data transaksi penjualan untuk hari ini
	db.Where("DATE(transaction_date) = ?", time.Now().Format("2006-01-02")).Find(&transaksiPenjualan)

	// Ambil data transaksi penjualan untuk minggu ini
	db.Where("WEEK(transaction_date) = ? AND YEAR(transaction_date) = ?", time.Now().Weekday(), time.Now().Year()).Find(&transaksiPenjualan)

	// Ambil data transaksi penjualan untuk bulan ini
	db.Where("MONTH(transaction_date) = ? AND YEAR(transaction_date) = ?", time.Now().Month(), time.Now().Year()).Find(&transaksiPenjualan)

	// Ambil data transaksi penjualan untuk tahun ini
	db.Where("YEAR(transaction_date) = ?", time.Now().Year()).Find(&transaksiPenjualan)

	return helper.ResponseBasic(c, 200, "Hallo world!")
}
