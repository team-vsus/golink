package models

type Channel struct {
	ID          uint `gorm:"primarykey"`
	Name        string
	CandidateID uint
	RecruiterID uint
	Messages    []Message
}
