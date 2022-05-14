package models

import (
	"gorm.io/gorm"
	"time"
	"github.com/team-vsus/golink/models"
)

type Interview struct {
	gorm.Model
	ApplicationID uint
	From time.Time
	Till time.Time
}