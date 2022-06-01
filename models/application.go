package models

import (
	"time"
)

type Application struct {
	ID         uint `gorm:"primarykey"`
	UserID     uint
	JobAdID    uint
	CreatedAt  time.Time
	Pinned     bool
	Documents  []Document
	Interviews []Interview
}
