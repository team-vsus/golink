package models

type Document struct {
	ID            uint `gorm:"primarykey"`
	Name          string
	Size          int
	ApplicationID uint
}
