package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
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

func GetAllJobAds (c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var jobAds []models.JobAd
	db.Find(&jobAds)

	c.JSON(200, jobAds)
}

// SUS
func GetJobAd (c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	var jobAd models.JobAd
	db.Find(&jobAd, "id = ?", id)

	c.JSON(200, jobAd)
}

func GetJobAdByCompany (c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	companyId := c.Param("companyId")

	var jobAds []models.JobAd
	db.Find(&jobAds, "company_id = ?", companyId)

	c.JSON(200, jobAds)
}

type createReq struct {
	Name string `json:"name"`
	Description string `json:"description"`
	Salary float64 `json:"salary"`
	CompanyId uint `json:"company_id"`
}

func (r createReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.Salary, validation.Required),
		validation.Field(&r.CompanyId, validation.Required),
	)
}

func CreateJobAd (c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req createReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var company models.Company
	result := db.Find(&company, "id = ?", req.CompanyId)
	if result.RowsAffected == 0 {
		c.JSON(400, "Company does not exist")
		return
	}

	newJobAd := &models.JobAd{
		Name: req.Name,
		Description: req.Description,
		Salary: req.Salary,
		CompanyId: req.CompanyId,
	}
	db.Create(newJobAd)

	c.JSON(200, "Successfully created new jobad")
}

func DeleteJobAd (c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	var jobAd models.JobAd
	result := db.Find(&jobAd, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(400, "JobAd does not exist")
		return
	}

	db.Delete(&jobAd)

	c.JSON(200, "Successfully deleted jobad")
}