package models

import "time"

type JobAd struct {
	ID           uint `gorm:"primarykey"`
	Name         string
	Description  string
	Salary       float64
	CompanyID    uint
	Open         bool
	Country      string
	City         string
	CreatedAt    time.Time
	Applications []Application
}
