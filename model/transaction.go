package model

import (
	"time"

	guuid "github.com/google/uuid"
)

type Transaction struct {
	ID              uint       `gorm:"primaryKey; autoIncrement" json:"-"`
	Price           float64    `json:"price"`
	Quantity        uint       `json:"qty"`
	Total           float64    `json:"-"`
	ReferProduct    guuid.UUID `json:"-"`
	TransactionDate time.Time  `gorm:"autoCreateTime" json:"-"`
}
