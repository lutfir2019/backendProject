package model

import (
	"time"

	guuid "github.com/google/uuid"
)

type Product struct {
	ID          int64         `gorm:"autoIncrement" json:"-"`
	PID         guuid.UUID    `gorm:"primaryKey; unique" json:"-"`
	Pnm         string        `json:"pnm"`
	Pcd         string        `json:"pcd"`
	Qty         int64         `json:"qty"`
	Price       int64         `json:"price"`
	Catcd       string        `json:"catcd"`
	Catnm       string        `json:"catnm"`
	Spcd        string        `json:"spcd"`
	Spnm        string        `json:"spnm"`
	Crby        string        `json:"crby"` // creat buy
	ShopRefer   guuid.UUID    `json:"-"`
	Transaction []Transaction `gorm:"foreignKey:ProductRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"-" `
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"-"`
}
