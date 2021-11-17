package db

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/team-vsus/golink/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func CreateConnection() *gorm.DB {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Europe/Vienna",
			os.Getenv("DBHOST"),
			os.Getenv("DBUSER"),
			os.Getenv("DBPW"),
			os.Getenv("DBNAME")),
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{})

	return db
}

func Inject(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
