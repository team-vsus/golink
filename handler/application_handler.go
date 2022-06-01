package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"
	"github.com/team-vsus/golink/models"
	"github.com/team-vsus/golink/utils"
)

func GetAllApplications(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var applications []models.Application
	db.Find(&applications)

	c.JSON(200, applications)
}

func GetApplicationByJobAd(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := jwt.MapClaims(c.MustGet("user").(jwt.MapClaims))["id"].(uint)
	applications := []models.Application{}
	db.Find(&applications, "job_ad_id = ?", id)

	c.JSON(200, applications)
}

func GetAllApplicationByUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")
	applications := []models.Application{}
	db.Find(&applications, "user_id = ?", id)

	c.JSON(200, applications)
}

func GetAllApplicationByMe(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.MustGet("user").(jwt.MapClaims)["id"].(uint)
	applications := []models.Application{}
	db.Find(&applications, "user_id = ?", id)

	c.JSON(200, applications)
}

type createApplicationReq struct {
	JobAdID   uint                `json:"job_ad_id"`
	Documents []createReqDocument `json:"documents"`
}

func (r createApplicationReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.JobAdID, validation.Required),
	)
}

// SUS docs
func CreateApplication(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req createApplicationReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var jobad models.JobAd
	result := db.Find(&jobad, "id = ?", req.JobAdID)
	if result.RowsAffected == 0 {
		c.JSON(400, "JobAd not found")
		return
	}

	for _, document := range req.Documents {
		newDocument := &models.Document{
			Name:          document.Name,
			Size:          document.Size,
			ApplicationID: document.ApplicationId,
		}
		db.Create(&newDocument)
	}

	application := &models.Application{
		JobAdID:   req.JobAdID,
		CreatedAt: time.Now(),
		Pinned:    false,
		UserID:    c.MustGet("user").(jwt.MapClaims)["id"].(uint),
	}
	db.Create(&application)

	c.JSON(200, application)
}

type deleteApplicationReq struct {
	ID uint `json:"id"`
}

func (r deleteApplicationReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.ID, validation.Required),
	)
}

// sus id param
func DeleteApplication(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req deleteApplicationReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var application models.Application
	result := db.Find(&application, "id = ?", req.ID)
	if result.RowsAffected == 0 {
		c.JSON(400, "Application not found")
		return
	}

	db.Delete(&application)

	c.JSON(200, "Successfully deleted application")
}

func DeleteApplicationbyJobAd(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")
	var application models.Application
	result := db.Find(&application, "job_ad_id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(400, "Application not found")
		return
	}

	db.Delete(&application)

	c.JSON(200, "Successfully deleted applications")
}
