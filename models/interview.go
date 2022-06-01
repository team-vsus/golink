package models

import (
	"time"
)

type Interview struct {
	ID            uint `gorm:"primarykey"`
	ApplicationID uint
	Date          time.Time
}
