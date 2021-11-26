package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string
	Username string
	Password string
	Locked   bool
	Verified bool
}
