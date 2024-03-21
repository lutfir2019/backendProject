package model

import (
	"time"

	guuid "github.com/google/uuid"
)

type Shop struct {
	ID          int64         `gorm:"autoIncrement" json:"-"`
	SID         guuid.UUID    `gorm:"primaryKey; unique" json:"-"`
	Spnm        string        `gorm:"unique" json:"spnm"`
	Spcd        string        `gorm:"unique" json:"spcd"`
	Almt        string        `json:"almt"`
	Product     []Product     `gorm:"foreignKey:ShopRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	User        []User        `gorm:"foreignKey:ShopRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Transaction []Transaction `gorm:"foreignKey:ShopRefer;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"-"`
}
