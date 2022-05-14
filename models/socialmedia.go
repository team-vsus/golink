package models

import (
	"gorm.io/gorm"
	"github.com/team-vsus/golink/models"
)

type SocialMedia struct {
	gorm.Model
	Link string
	CompanyID uint
}