package models

import (
	"gorm.io/gorm"
	"time"
	"github.com/team-vsus/golink/models"
)

type Application struct {
	gorm.Model
	UserID uint
	JobAdID uint 
	createdAt time.Time
	Documents []models.Document
	Interviews []models.Interview
}