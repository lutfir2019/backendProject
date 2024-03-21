package model

import (
	"time"

	guuid "github.com/google/uuid"
)

type Transaction struct {
	ID              uint       `gorm:"primaryKey; autoIncrement" json:"-"`
	Price           int64      `json:"-"`
	Quantity        int64      `json:"-"`
	Total           int64      `json:"-"`
	ProductRefer    guuid.UUID `json:"-"`
	ShopRefer       guuid.UUID `json:"-"`
	TransactionDate time.Time  `gorm:"autoCreateTime" json:"-"`
}
