package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/team-vsus/golink/models"
	"github.com/team-vsus/golink/utils"
	"gorm.io/gorm"
)

func GetAllSocialMedias(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var socialmedia []models.SocialMedia
	db.Find(&socialmedia)
	c.JSON(200, socialmedia)
}

func GetSocialMediaByCompanyId(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	companyid := c.Param("companyid")

	var socialmedia models.SocialMedia

	err := db.First(&socialmedia, "company_id = ?", companyid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, "SocialMedia not found")
		return
	}

	c.JSON(200, socialmedia)
}

type createReqSocialMedia struct {
	Link      string `json:"link"`
	CompanyId uint   `json:"company_id"`
}

func (r createReqSocialMedia) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Link, validation.Required),
		validation.Field(&r.CompanyId, validation.Required),
	)
}

func createSocialMedia(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req createReqSocialMedia
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var socialmedia models.SocialMedia
	result := db.Find(&socialmedia, "id = ?", req.CompanyId)
	if result.RowsAffected == 0 {
		c.JSON(400, "Company does not exist")
		return
	}

	newSocialMedia := &models.SocialMedia{
		Link:      req.Link,
		CompanyID: req.CompanyId,
	}

	db.Create(newSocialMedia)

	c.JSON(200, "Successfully created new SocialMedia")
}

func deleteSocialMedia(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var socialmedia models.SocialMedia
	result := db.Find(&socialmedia, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(400, "SocialMedia does not exist")
		return
	}

	db.Delete(&socialmedia)

	c.JSON(200, "Successfully deleted Social Media")
}

func deleteAllSocialMediasByCompanyId(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req createReqSocialMedia
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var socialmedia []models.SocialMedia

	result := db.Find(&socialmedia, "company_id = ?", req.CompanyId)
	if result.RowsAffected == 0 {
		c.JSON(400, "Company does not exist")
		return
	}

	db.Delete(&socialmedia)

	c.JSON(200, "Successfully deleted all social medias from company")
}
