package models

import (
	"gorm.io/gorm"
	"time"
	"github.com/team-vsus/golink/models"
)

type Message struct {
	gorm.Model
	ChannelID uint
	Content string
	SenderID uint
	createdAt time.Time
}