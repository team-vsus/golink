package handler

import (
	"math/rand"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/gorm"

	"github.com/team-vsus/golink/models"
	"github.com/team-vsus/golink/utils"
)

func GetAllCompanies(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var companies []models.Company
	db.Find(&companies)

	c.JSON(200, companies)
}

func GetMyCompany(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	db.First(&user, "id = ?", uint(c.MustGet("user").(jwt.MapClaims)["id"].(float64)))

	var company models.Company
	result := db.Find(&company, "id = ?", user.CompanyID)
	if result.RowsAffected == 0 {
		c.JSON(400, "Company does not exist")
		return
	}

	c.JSON(200, company)
}

func GetCompany(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	var company models.Company
	db.Find(&company, "id = ?", id)

	c.JSON(200, company)
}

type createCompanyReq struct {
	Name       string `json:"name"`
	UserId     int    `json:"userId"`
	WebsiteUrl string `json:"websiteUrl"`
	Address    string `json:"address"`
	Country    string `json:"country"`
}

func (r createCompanyReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
	)
}

// sus owner
func CreateCompany(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req createCompanyReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var company models.Company
	result := db.Find(&company, "name = ?", req.Name)
	if result.RowsAffected != 0 {
		c.JSON(400, "Company already exists")
		return
	}

	var id = uint(c.MustGet("user").(jwt.MapClaims)["id"].(float64))

	newCompany := &models.Company{
		Name:       req.Name,
		OwnerID:    uint(req.UserId),
		Country:    req.Country,
		WebsiteUrl: req.WebsiteUrl,
		Address:    req.Address,
	}
	db.Create(&newCompany)

	var user models.User
	db.First(&user, "id = ?", id)
	user.CompanyID = newCompany.ID
	db.Save(&user)

	c.JSON(200, "Successfully created new company")
}

type deleteCompanyReq struct {
	Id string `json:"id"`
}

func GenerateCompanyInvite(db *gorm.DB, company_id int) int {
	code := 10000000 + rand.Intn(99999999-10000000)

	db.Create(&models.CompanyInvite{
		Code:      code,
		CompanyId: company_id,
		ExpiresAt: time.Now().AddDate(0, 0, 7),
	})

	return code
}

func DeleteCompany(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var company models.Company
	db.Find(&company, "user_id = ?", uint(c.MustGet("user").(jwt.MapClaims)["id"].(float64)))

	db.Delete(&company)

	c.JSON(200, "Successfully deleted company")
}

func GetCompanyInvite(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var user models.User
	db.First(&user, "id = ?", uint(c.MustGet("user").(jwt.MapClaims)["id"].(float64)))

	var company models.Company
	result := db.Find(&company, "id = ?", user.CompanyID)
	if result.RowsAffected == 0 {
		c.JSON(400, "Company does not exist")
		return
	}

	code := GenerateCompanyInvite(db, int(user.CompanyID))

	c.JSON(200, code)
}
