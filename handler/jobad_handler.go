package handler

import (


	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"gorm.io/gorm"

	"github.com/team-vsus/golink/models"
	"github.com/team-vsus/golink/utils"
)

func GetAllJobAds(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var jobAds []models.JobAd
	if filter := c.Param("filter"); filter == "" {
		db.Find(&jobAds)
	} else {
		db.Where("open = ?", filter).Find(&jobAds)
	}

	c.JSON(200, jobAds)
}

// SUS
func GetJobAd(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	var jobAd models.JobAd
	db.Find(&jobAd, "id = ?", id)

	c.JSON(200, jobAd)
}


func GetJobAdSearch(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var jobAds []models.JobAd
	db.Where("description LIKE ?", "%"+c.Param("search")+"%").Find(&jobAds)

	c.JSON(200, jobAds)
}


func GetJobAdByCompany(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	companyId := c.Param("companyId")

	var jobAds []models.JobAd
	db.Find(&jobAds, "company_id = ?", companyId)

	c.JSON(200, jobAds)
}


type salaryReqJobAd struct {
	Lowersalary float64 `json:"lower_salary"`
	Uppersalary float64 `json:"upper_salary"`
}

func (r salaryReqJobAd) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Lowersalary, validation.Required),
		validation.Field(&r.Uppersalary, validation.Required),
	)
}

// sus
func GetJobAdBySalary(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req salaryReqJobAd
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var jobAds []models.JobAd
	db.Where("salary >= ? AND salary <= ?", req.Lowersalary, req.Uppersalary).Find(&jobAds)

	c.JSON(200, jobAds)
}


type createJobAdReq struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Salary      float64 `json:"salary"`
	CompanyId   uint    `json:"company_id"`
	Location    string  `json:"location"`
}

func (r createJobAdReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Description, validation.Required),
		// validation.Field(&r.Salary, validation.Required),
	)
}


func CreateJobAd(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req createJobAdReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var user models.User
	db.First(&user, "id = ?", c.MustGet("user").(jwt.MapClaims)["id"].(uint))

	var company models.Company
	result := db.Find(&company, "id = ?", user.CompanyID)
	if result.RowsAffected == 0 {
		c.JSON(400, "Company does not exist")
		return
	}

	newJobAd := &models.JobAd{
		Name:        req.Name,
		Description: req.Description,
		Salary:      req.Salary,
		CompanyID:   req.CompanyId,
		Open:        true,
		Location:    req.Location,
		CompanyID:   company.ID,
	}
	db.Create(newJobAd)

	c.JSON(200, "Successfully created new jobad")
}

func DeleteJobAd(c *gin.Context) {
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
