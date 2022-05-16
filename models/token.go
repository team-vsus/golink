package models

import (
	"time"
)

type Token struct {
	ID        uint `gorm:"primarykey"`
	Token     string
	UserId    int
	ExpiresAt time.Time
}
