package models

import (
	"gorm.io/gorm"
	"github.com/team-vsus/golink/models"
)

type User struct {
	gorm.Model
	Email     string
	Firstname string
	Lastname  string
	Password  string
	Locked    bool
	Verified  bool
	Role 	  int
	CompanyID uint
	Applications []models.Application
	ChannelsCandidate []models.Channel `gorm:"foreignKey:CandidateID"`
	ChannelsRecruiter []models.Channel `gorm:"foreignKey:RecruiterID"`
	Messages []models.Message `gorm:"foreignKey:SenderID"`
}
