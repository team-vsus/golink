package models

import	(
	"gorm.io/gorm"
	"github.com/team-vsus/golink/models"
)

type Channel struct {
	gorm.Model
	Name string
	CandidateID uint
	RecruiterID uint
	Messages []models.Message
}