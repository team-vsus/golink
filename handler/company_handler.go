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

func GetAllCompanies(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var companies []models.Company
	db.Find(&companies)

	c.JSON(200, companies)
} 

func GetCompany (c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	var company models.Company
	db.Find(&company, "id = ?", id)

	c.JSON(200, company)
}

type createReq struct {
	Name string `json:"name"`
}

func (r createReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
	)
}

func CreateCompany(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req createReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var company models.Company
	result := db.Find(&company, "name = ?", req.Name)
	if result.RowsAffected != 0 {
		c.JSON(400, "Company already exists")
		return
	}

	newCompany := &models.Company{
		Name: req.Name,
	}
	db.Create(&newCompany)

	c.JSON(200, "Successfully created new company")
}

type deleteReq struct {
	Id string `json:"id"`
}
// sus
func DeleteCompany (c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	var company models.Company
	db.Find(&company, "id = ?", id)

	db.Delete(&company)

	c.JSON(200, "Successfully deleted company")
}

