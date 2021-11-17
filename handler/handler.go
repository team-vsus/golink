package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/team-vsus/golink/db"
)

func InitHandler() *gin.Engine {
	r := gin.Default()

	conn := db.CreateConnection()
	r.Use(db.Inject(conn))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "Pong")
	})

	// Example
	/* ug := r.Group("/api/v1/users")
	   ug.GET("", GetAllUsers)
	   ug.GET(":id", GetUser) **/

	return r
}
