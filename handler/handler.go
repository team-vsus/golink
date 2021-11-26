package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/team-vsus/golink/db"
	"github.com/team-vsus/golink/utils"
)

func InitHandler() *gin.Engine {
	r := gin.Default()

	conn := db.CreateConnection()
	r.Use(db.Inject(conn))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "Pong")
	})

	ag := r.Group("/api/v1/auth")
	ag.POST("/register", Register)
	ag.POST("/login", Login)
	ag.POST("/verify", Verify)
	ag.POST("/forgot/password/new", ForgotPasswordNew)
	ag.POST("/forgot/password/", ForgotPassword)

	ug := r.Group("/api/v1/users")
	ug.Use(utils.VerifyToken)
	ug.GET("", GetAllUsers)
	ug.GET(":id", GetUser)
	ug.GET("/me", GetMe)

	return r
}
