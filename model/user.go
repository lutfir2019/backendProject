package model

import (
	"time"

	guuid "github.com/google/uuid"
)

type User struct {
	ID        int64      `gorm:"autoIncrement" json:"-"`
	UID       guuid.UUID `gorm:"primaryKey; unique" json:"-"`
	Nam       string     `json:"nam"`
	Unm       string     `gorm:"unique" json:"unm"`
	Pass      string     `json:"-"`
	Rlcd      string     `json:"rlcd"`
	Rlnm      string     `json:"rlnm"`
	Almt      string     `json:"almt"`
	Gdr       string     `json:"gdr"`
	Pn        string     `json:"pn"`
	Spcd      string     `json:"spcd"`
	Spnm      string     `json:"spnm"`
	ShopRefer guuid.UUID `json:"-"`
	Sessions  []Session  `gorm:"foreignKey:UserRefer; constraint:OnUpdate:CASCADE, OnDelete:CASCADE;" json:"-"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"-" `
	UpdatedAt time.Time  `gorm:"autoUpdateTime:milli" json:"-"`
}
