package models

type SocialMedia struct {
	ID        uint `gorm:"primarykey"`
	Link      string
	CompanyID uint
}
