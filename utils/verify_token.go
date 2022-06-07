package utils

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func VerifyToken(c *gin.Context) {
	//tokenString := c.GetHeader("authorization")
	tokenString, err := c.Cookie("token")

	if len(tokenString) == 0 || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"failed": true,
		})
		c.Abort()
		return
	}

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		if err := token.Claims.Valid(); err != nil {
			return nil, err
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	fmt.Println("Claims", claims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"failed": true,
		})
		c.Abort()
		return
	}
	c.Set("user", claims)
	c.Next()
}
