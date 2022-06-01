package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/team-vsus/golink/models"
	"github.com/team-vsus/golink/utils"
)

func GetAllApplications (c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var applications []models.Application
	db.Find(&applications)

	c.JSON(200, applications)
}

func GetApplicationByJobAd (c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")
	applications := []models.Application{}
	db.Find(&applications, "job_ad_id = ?", id)

	c.JSON(200, applications)
}

type createReq struct {
	JobAdID uint `json:"job_ad_id"`
	Documents []models.Document `json:"documents"`
}

func (r createReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.JobAdID, validation.Required),
	)
}


func CreateApplication (c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)


	var req createReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}


	var jobad models.JobAd
	result := db.Find(&jobad, "id = ?", req.JobAdID)
	if result.RowsAffected == 0 {
		c.JSON(400, "JobAd not found")
		return
	}

	newApplication := &models.Application{
		JobAdID: req.JobAdID,
		Documents: req.Documents,
		CreatedAt: time.Now(),
		Pinned: false,
		UserID: c.MustGet("user").(jwt.MapClaims)["id"].(uint),
	}
	db.Create(newApplication)

	uploadDocument(db,req.Documents)

	 c.JSON(200, newApplication)
}

type deleteReq struct {
	ID uint `json:"id"`
}

func DeleteApplication (c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req deleteReq
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

func uploadDocument(db *gorm.DB, documents [] models.Document){
		form, err := c.MultipartForm()
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}
		file := form.File["files"]

		for _, file := range documents {
			filename := filepath.Base(file.Name)
			if err := c.SaveUploadedFile(file, filename); err != nil {
				c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
				return
			}
		}
}