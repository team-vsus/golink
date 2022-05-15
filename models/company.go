package models

import (
	"gorm.io/gorm"
	"github.com/team-vsus/golink/models"
)

type Company struct {
	gorm.Model
	Name string
	UserID uint
	JobAds []models.JobAd
	SocialMedias []models.SocialMedia
	Users []models.User
}