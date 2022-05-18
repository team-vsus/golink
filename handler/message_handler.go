package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/team-vsus/golink/models"
	"github.com/team-vsus/golink/utils"
	"gorm.io/gorm"
)

func GetAllMessages(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var message []models.Message
	db.Find(&message)
	c.JSON(200, message)
}

func GetMessageByChannelId(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	channelid := c.Param("channelid")

	var message models.Message

	err := db.First(&message, "channel_id = ?", channelid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, "Message not found")
		return
	}

	c.JSON(200, message)
}

type createReqMessage struct {
	ChannelId uint   `json:"channel_id"`
	Content   string `json:"content"`
	SenderId  uint   `json:"sender_id"`
}

func (r createReqMessage) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ChannelId, validation.Required),
		validation.Field(&r.Content, validation.Required),
		validation.Field(&r.SenderId, validation.Required),
	)
}

// sus sender
func createMessage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req createReqMessage
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var message models.Message
	result := db.Find(&message, "id = ?", req.ChannelId)
	if result.RowsAffected == 0 {
		c.JSON(400, "Channel does not exist")
		return
	}

	newMessage := &models.Message{
		ChannelID: req.ChannelId,
		Content:   req.Content,
		SenderID:  req.SenderId,
	}

	db.Create(newMessage)

	c.JSON(200, "Successfully created new Message")
}

func deleteMessage(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var message models.Message
	result := db.Find(&message, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(400, "Message does not exist")
		return
	}

	db.Delete(&message)

	c.JSON(200, "Successfully deleted Message")
}

type deleteReqMessage struct {
	ChannelId uint `json:"channel_id"`
}

func (r deleteReqMessage) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ChannelId, validation.Required),
	)
}

func deleteAllMessagesByChannelId(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req deleteReqMessage
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var message []models.Message

	result := db.Find(&message, "channel_id = ?", req.ChannelId)
	if result.RowsAffected == 0 {
		c.JSON(400, "Channel does not exist")
		return
	}

	db.Delete(&message)

	c.JSON(200, "Successfully deleted all messages from channel")

}
