package model

import (
	guuid "github.com/google/uuid"
)

type User struct {
	ID        int64      `gorm:"primaryKey" json:"-"`
	UID       guuid.UUID `json:"-"`
	Nam       string     `json:"nam"`
	Unm       string     `json:"unm"`
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
	CreatedAt int64      `gorm:"autoCreateTime" json:"-" `
	UpdatedAt int64      `gorm:"autoUpdateTime:milli" json:"-"`
}
