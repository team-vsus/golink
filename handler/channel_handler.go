package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/gorm"

	"github.com/team-vsus/golink/models"
	"github.com/team-vsus/golink/utils"
)

func GetAllChannels(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var channels []models.Channel
	db.Find(&channels)

	c.JSON(200, channels)
}

func GetChannel(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	var channel models.Channel
	db.Find(&channel, "id = ?", id)

	c.JSON(200, channel)
}

func GetChannelByUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	userID := c.MustGet("user").(jwt.MapClaims)["id"].(uint)

	var channel models.Channel
	db.Find(&channel, "candidate_id = ? or recruiter_id = ?", userID, userID)

	c.JSON(200, channel)
}

type createChannelReq struct {
	Name        string `json:"name"`
	CandidateID uint   `json:"candidate_id"`
	RecruiterID uint   `json:"recruiter_id"`
}

func (r createChannelReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.CandidateID, validation.Required),
		validation.Field(&r.RecruiterID, validation.Required),
	)
}

func CreateChannel(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req createChannelReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	newChannel := &models.Channel{
		Name:        req.Name,
		CandidateID: req.CandidateID,
		RecruiterID: req.RecruiterID,
	}

	db.Create(newChannel)

	c.JSON(200, newChannel)
}

func DeleteChannel(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	var channel models.Channel
	result := db.Find(&channel, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(400, "Channel does not exist")
		return
	}

	db.Delete(&channel)

	c.JSON(200, "Successfully deleted channel")
}
