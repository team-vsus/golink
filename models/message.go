package models

import (
	"time"
)

type Message struct {
	ID        uint `gorm:"primarykey"`
	ChannelID uint
	Content   string
	SenderID  uint
	createdAt time.Time
}
