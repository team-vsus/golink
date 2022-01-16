package handler

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/team-vsus/golink/models"
	"gorm.io/gorm"
)

func GetAllUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var users []models.User
	db.Find(&users)
	c.JSON(200, users)
}

func GetUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var user models.User

	// Check if returns RecordNotFound error
	err := db.First(&user, "Id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, "User not found")
		return
	}

	c.JSON(200, user)
}

func GetMe(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	jwtPayload := c.MustGet("user").(jwt.MapClaims)
	userId := jwtPayload["id"]

	var user models.User
	err := db.First(&user, "id = ?", userId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, gin.H{
			"failed": true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"firstname": user.Firstname,
		"lastname":  user.Lastname,
	})
}
