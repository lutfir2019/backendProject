package model

import (
	"time"

	guuid "github.com/google/uuid"
)

type Session struct {
	Sessionid guuid.UUID `gorm:"primaryKey" json:"sessionid"`
	Expires   time.Time  `json:"-"`
	UserRefer guuid.UUID `json:"-"`
	Token     string     `json:"token"` // JWT token for client side to access the session data
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"-" `
}
