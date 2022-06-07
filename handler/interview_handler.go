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
	ApplicationId uint   `json:"application_id"`
	Date          string `json:"date"`
}

func (r creatReqInterview) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ApplicationId, validation.Required),
		validation.Field(&r.Date, validation.Required),
	)
}

func createInterview(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req creatReqInterview
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var application models.Application
	result := db.Find(&application, "id = ?", req.ApplicationId)
	if result.RowsAffected == 0 {
		c.JSON(400, "Application does not exist")
		return
	}

	var date, err = time.Parse("2006-01-02 03:04:05", req.Date)
	if err != nil {
		c.JSON(400, "Invalid date format")
		return
	}

	newInterview := &models.Interview{
		ApplicationID: req.ApplicationId,
		Date:          date,
	}

	db.Create(newInterview)

	c.JSON(200, "Successfully created new Interview")
}

type deleteReqInterview struct {
	ApplicationId uint `json:"application_id"`
}

func (r deleteReqInterview) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ApplicationId, validation.Required),
	)
}

func deleteAllInterviewsByApplicationId(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req deleteReqInterview
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

type patchReqInterview struct {
	ID   uint      `json:"id"`
	Date time.Time `json:"date"`
}

func (r patchReqInterview) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
		validation.Field(&r.Date, validation.Required),
	)
}

func updateInterviewDate(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req patchReqInterview
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var interview models.Interview

	result := db.Find(&interview, "id = ?", req.ID)
	if result.RowsAffected == 0 {
		c.JSON(400, "Interview does not exist")
		return
	}

	interview.Date = req.Date

	db.Save(&result)

	c.JSON(200, "Successfully updated date from Interview")
}
