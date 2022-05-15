package models

import (
	"gorm.io/gorm"
	"github.com/team-vsus/golink/models"
)

type JobAd struct {
	gorm.Model
	Name string
	Description string
	Salary float64
	CompanyID uint
	Applications []models.Application
}