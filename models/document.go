package models

import (
	"gorm.io/gorm"
	"github.com/team-vsus/golink/models"
)

type Document struct {
	gorm.Model
	Name string
	Size int
	ApplicationID uint
}