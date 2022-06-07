package models

type Company struct {
	ID           uint `gorm:"primarykey"`
	Name         string
	OwnerID      uint
	WebsiteUrl   string
	Address      string
	Country      string
	JobAds       []JobAd
	SocialMedias []SocialMedia
	Users        []User `gorm:"constraint:OnUpdate:SET NULL,OnDelete:SET NULL;"`
}
