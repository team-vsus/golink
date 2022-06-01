package models

type Company struct {
	ID           uint `gorm:"primarykey"`
	Name         string
	UserID       uint
	JobAds       []JobAd
	SocialMedias []SocialMedia
	Users        []User
}
