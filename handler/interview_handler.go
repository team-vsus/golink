package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/team-vsus/golink/models"
	"github.com/team-vsus/golink/utils"
	"gorm.io/gorm"
)

func GetAllInterviews(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var interview []models.Interview
	db.Find(&interview)
	c.JSON(200, interview)
}

func GetInterviewByApplicationId(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	applicationid := c.Param("applicationid")

	var interview models.Interview

	err := db.First(&interview, "application_id", applicationid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, "Interview not found")
		return
	}

	c.JSON(200, interview)
}

type creatReqInterview struct {
	ApplicationId uint      `json:"application_id"`
	From          time.Time `json:"from"`
	Till          time.Time `json:"till"`
}

func (r creatReqInterview) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ApplicationId, validation.Required),
		validation.Field(&r.From, validation.Required),
		validation.Field(&r.Till, validation.Required),
	)
}

func createInterview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req creatReqInterview
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var interview models.Interview
	result := db.Find(&interview, "id = ?", req.ApplicationId)
	if result.RowsAffected == 0 {
		c.JSON(400, "Application does not exist")
		return
	}

	newInterview := &models.Interview{
		ApplicationID: req.ApplicationId,
		From:          req.From,
		Till:          req.Till,
	}

	db.Create(newInterview)

	c.JSON(200, "Successfully created new Interview")
}

func deleteAllInterviewsByApplicationId(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req creatReqInterview
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var interview []models.Interview

	result := db.Find(&interview, "application_id = ?", req.ApplicationId)
	if result.RowsAffected == 0 {
		c.JSON(400, "Application does not exist")
		return
	}

	db.Delete(&interview)

	c.JSON(200, "Successfully deleted all interviews from  application")
}
