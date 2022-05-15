package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/team-vsus/golink/models"
	"github.com/team-vsus/golink/utils"
	"gorm.io/gorm"
)

func GetAllDocuments(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var document []models.Document
	db.Find(&document)
	c.JSON(200, document)
}

func GetDocumentByApplicationId(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	applicationid := c.Param("applicationid")

	var document models.Document

	err := db.First(&document, "application_id = ?", applicationid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, "Document not found")
		return
	}

	c.JSON(200, document)
}

type createReqDocument struct {
	Name          string `json:"name"`
	Size          int    `json:"size"`
	ApplicationId uint   `json:"application_id"`
}

func (r createReqDocument) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Size, validation.Required),
		validation.Field(&r.ApplicationId, validation.Required),
	)

}

func createDocument(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req createReqDocument
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var document models.Document
	result := db.Find(&document, "id = ?", req.ApplicationId)
	if result.RowsAffected == 0 {
		c.JSON(400, "Application does not exist")
		return
	}

	newDocument := &models.Document{
		Name:          req.Name,
		Size:          req.Size,
		ApplicationID: req.ApplicationId,
	}

	db.Create(newDocument)

	c.JSON(200, "Successfully created new Document")
}

func deleteDocument(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var document models.Document
	result := db.Find(&document, "id = ?", id)
	if result.RowsAffected == 0 {
		c.JSON(400, "Document does not exist")
		return
	}

	db.Delete(&document)

	c.JSON(200, "Successfully deleted Document")
}

func deleteAllDocumentByApplicationId(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req createReqDocument
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var document []models.Document

	result := db.Find(&document, "application_id = ?", req.ApplicationId)
	if result.RowsAffected == 0 {
		c.JSON(400, "Application does not exist")
		return
	}

	db.Delete(document)

	c.JSON(200, "Successfully deleted all documents")
}
