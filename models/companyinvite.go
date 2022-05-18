package models

import (
	"time"
)

type CompanyInvite struct {
	ID        uint `gorm:"primarykey"`
	Code      int
	CompanyId int
	ExpiresAt time.Time
}
