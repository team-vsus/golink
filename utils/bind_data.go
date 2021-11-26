package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Request interface {
	Validate() error
}

func BindData(c *gin.Context, req Request) bool {
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	if err := req.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return false
	}

	return true
}
