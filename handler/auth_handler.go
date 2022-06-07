package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"github.com/team-vsus/golink/models"
	"github.com/team-vsus/golink/utils"
)

type registerReq struct {
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
	Applicant bool   `json:"applicant"`
}

func (r registerReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required),
		validation.Field(&r.Firstname, validation.Required, validation.Length(2, 20)),
		validation.Field(&r.Lastname, validation.Required, validation.Length(2, 20)),
		validation.Field(&r.Password, validation.Required, validation.Length(5, 30)),
	)
}

func Register(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// validate input
	var req registerReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	// check if the user doesn't already exists
	var user models.User
	result := db.Find(&user, "email = ?", req.Email)
	if result.RowsAffected != 0 {
		c.JSON(400, "User already exists")
		return
	}

	// hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		c.JSON(500, "Error while creating new user")
		return
	}

	// insert to database
	newUser := &models.User{
		Email:     req.Email,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Password:  string(hash),
		Applicant: req.Applicant,
	}
	db.Create(&newUser)

	// generate new token
	token := utils.GenerateNewToken(6)
	db.Create(&models.Token{
		Token:     token,
		UserId:    int(newUser.ID),
		ExpiresAt: time.Now().AddDate(0, 0, 7),
	})

	// send confirmation email
	e := email.NewEmail()
	e.From = "muazahmed766@outlook.com"
	e.To = []string{"muazahmed019@gmail.com"}
	e.Subject = "Veriy your email"
	e.Text = []byte(token)
	go utils.SendEmail(e)
	fmt.Println("Sending email")

	c.JSON(200, "Successfully created new user")

}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r loginReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required),
		validation.Field(&r.Password, validation.Required),
	)
}

func Login(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	// validate input
	var req loginReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}
	// check if user exists
	var user models.User
	err := db.First(&user, "email = ?", req.Email).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, "User not found")
		return
	}
	fmt.Println(user)
	// check if user is locked
	if user.Locked {
		c.JSON(400, "User locked!")
		return
	}
	// check if user is verified
	if !user.Verified {
		c.JSON(400, "User not verified!")
		return
	}

	// check if password is correct
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(404, "Password not correct")
		return
	}

	// generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat":       time.Now().Unix(),
		"id":        user.ID,
		"firstname": user.Firstname,
		"lastname":  user.Lastname,
		"applicant": user.Applicant,
	})

	tokenStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(400, "Something went wrong!")
		return
	}

	c.SetCookie("token", tokenStr, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"email":     user.Email,
		"firstname": user.Firstname,
		"lastname":  user.Lastname,
		"applicant": user.Applicant,
	})
}

type verifyReq struct {
	Code string `json:"code"`
}

func (r verifyReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Code, validation.Required),
	)
}

func Verify(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req verifyReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var token models.Token
	err := db.First(&token, "token = ?", req.Code).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(400, "Code expired")
		return
	}

	if token.ExpiresAt.After(token.ExpiresAt.AddDate(0, 0, 7)) {
		c.JSON(400, "Code expired")
		return
	}

	db.Model(&models.User{}).Where("id = ?", token.UserId).Update("verified", true)

	c.JSON(200, gin.H{
		"id": token.UserId,
	})
}

type forgotPasswordReq struct {
	Email string `json:"email"`
}

func (r forgotPasswordReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required),
	)
}

func ForgotPassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req forgotPasswordReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var user models.User
	err := db.First(&user, "email = ?", req.Email).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, "User not found!")
		return
	}

	token := uuid.New()
	db.Create(&models.Token{
		Token:     token.String(),
		UserId:    int(user.ID),
		ExpiresAt: time.Now().AddDate(0, 0, 7),
	})
	link := fmt.Sprintf("%s/auth/reset-pw/%s", os.Getenv("FRONTEND_HOST"), token)
	println(link)

	//utils.SendEmail("pikayuhno@gmail.com", []string{"muazahmed019@gmail.com"}, []byte(link))

	c.JSON(200, "Successfully sent link for changing password!")
}

type resetPasswordReq struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

func (r resetPasswordReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Token, validation.Required),
		validation.Field(&r.NewPassword, validation.Required),
	)
}

func ResetPassword(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var req resetPasswordReq
	if ok := utils.BindData(c, &req); !ok {
		return
	}

	var token models.Token
	db.First(&token, "token = ?", req.Token)

	println("Token", token.Token)
	println("Token", token.UserId)

	var user models.User
	err := db.First(&user, "id = ?", token.UserId).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(404, "User not found!")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), 10)
	if err != nil {
		c.JSON(500, "Error while creating new user")
	}

	user.Password = string(hash)
	db.Model(&models.User{}).Where("id = ?", token.UserId).Update("password", string(hash))

	c.JSON(200, "Successfully changed password!")
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", 0, "/", "localhost", false, true)
	c.JSON(http.StatusOK, "")
}
