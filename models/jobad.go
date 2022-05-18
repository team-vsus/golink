package models

type JobAd struct {
	ID           uint `gorm:"primarykey"`
	Name         string
	Description  string
	Salary       float64
	CompanyID    uint
	Open         bool
	Location     string
	Applications []Application
}
