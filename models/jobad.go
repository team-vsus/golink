package models

type JobAd struct {
	ID           uint `gorm:"primarykey"`
	Name         string
	Description  string
	Salary       float64
	CompanyID    uint
	Applications []Application
}
