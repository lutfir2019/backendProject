package model

import guuid "github.com/google/uuid"

type Shop struct {
	ID        int64      `gorm:"autoIncrement" json:"-"`
	SID       guuid.UUID `gorm:"primaryKey" json:"-"`
	Spnm      string     `json:"spnm"`
	Spcd      string     `json:"spcd"`
	Almt      string     `json:"almt"`
	Product   []Product  `gorm:"foreignKey:ShopRefer; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"-"`
	User      []User     `gorm:"foreignKey:ShopRefer; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"-"`
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64      `gorm:"autoUpdateTime" json:"-"`
}
