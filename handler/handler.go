package handler

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/team-vsus/golink/db"
	"github.com/team-vsus/golink/utils"
)

func InitHandler() *gin.Engine {
	r := gin.Default()

	conn := db.CreateConnection()
	//r.Use(cors.Default())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	}))
	r.Use(db.Inject(conn))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, "Pong")
	})

	ag := r.Group("/api/v1/auth")
	ag.POST("/register", Register)
	ag.POST("/login", Login)
	ag.POST("/logout", Logout)
	ag.POST("/verify", Verify)
	ag.POST("/forgot/password/new", ForgotPasswordNew)
	ag.POST("/forgot/password/", ForgotPassword)

	ug := r.Group("/api/v1/users")
	ug.Use(utils.VerifyToken)
	ug.GET("", GetAllUsers)
	ug.GET(":id", GetUser)
	ug.GET("/me", GetMe)

	cg := r.Group("/api/v1/companys")
	cg.Use(utils.VerifyToken)
	cg.GET("", GetAllCompanies)
	cg.GET(":id", GetCompany)
	cg.GET("/invite/", GetCompanyInvite)
	cg.POST("", CreateCompany)
	cg.DELETE("", DeleteCompany)

	apg := r.Group("/api/v1/applications")
	apg.Use(utils.VerifyToken)
	apg.GET("", GetAllApplicationByMe)
	apg.GET("/user/:id", GetAllApplicationByUser)
	apg.GET("/all/", GetAllApplications)
	apg.GET("/job/:id", GetApplicationByJobAd)
	apg.POST("", CreateApplication)
	apg.DELETE("/:id", DeleteApplication)
	apg.DELETE("/job/:id", DeleteApplicationbyJobAd)

	jg := r.Group("/api/v1/jobads")
	jg.Use(utils.VerifyToken)
	// getJobAdByMe
	jg.GET("", GetAllJobAds)
	jg.GET(":id", GetJobAd)
	jg.GET("/company/:id", GetJobAdByCompany)
	jg.GET("/search/:search", GetJobAdSearch)
	jg.GET("/salary/", GetJobAdBySalary)
	jg.POST("", CreateJobAd)
	jg.DELETE("/:id", DeleteJobAd)

	dg := r.Group("/api/v1/documents")
	dg.Use(utils.VerifyToken)
	dg.GET("", GetAllDocuments)
	dg.GET(":applicationid", GetDocumentByApplicationId)
	dg.POST("", createDocument)
	dg.DELETE("/:id", deleteDocument)
	dg.DELETE("", deleteAllDocumentByApplicationId)

	ig := r.Group("/api/v1/interviews")
	ig.Use(utils.VerifyToken)
	ig.GET("", GetAllInterviews)
	ig.GET(":applicationid", GetInterviewByApplicationId)
	ig.POST("", createInterview)
	ig.PATCH("", updateInterviewDate)
	ig.DELETE("/", deleteAllInterviewsByApplicationId)

	mg := r.Group("/api/v1/messages")
	mg.Use(utils.VerifyToken)
	mg.GET("", GetAllMessages)
	mg.GET(":channelid", GetMessageByChannelId)
	mg.POST("", createMessage)
	mg.DELETE("/:id", deleteMessage)
	mg.DELETE("", deleteAllMessagesByChannelId)

	sg := r.Group("/api/v1/social")
	sg.Use(utils.VerifyToken)
	sg.GET("", GetAllSocialMedias)
	sg.GET(":companyid", GetSocialMediaByCompanyId)
	sg.POST("", createSocialMedia)
	sg.DELETE("/:id", deleteSocialMedia)
	sg.DELETE("", deleteAllSocialMediasByCompanyId)

	chg := r.Group("/api/v1/channels")
	chg.Use(utils.VerifyToken)
	chg.GET("", GetAllChannels)
	chg.GET("/user/", GetChannelByUser)
	chg.POST("", CreateChannel)
	chg.DELETE("/:id", DeleteChannel)

	return r
}
