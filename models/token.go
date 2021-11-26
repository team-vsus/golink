package models

import (
	"time"

	"gorm.io/gorm"
)

type Token struct {
	gorm.Model
	Token     string
	UserId    int
	ExpiresAt time.Time
}
