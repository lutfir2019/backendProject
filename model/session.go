package model

import (
	"time"

	guuid "github.com/google/uuid"
)

type Session struct {
	Sessionid guuid.UUID `gorm:"primaryKey; unique" json:"sessionid"`
	Expires   time.Time  `json:"-"`
	UserRefer guuid.UUID `json:"-"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"-" `
}
