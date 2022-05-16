package models

type User struct {
	ID        uint `gorm:"primarykey"`
	Email     string
	Firstname string
	Lastname  string
	Password  string
	Locked    bool
	Verified  bool
	Role      int
	CompanyID uint
	// Description       string
	Applications      []Application
	ChannelsCandidate []Channel `gorm:"foreignKey:CandidateID"`
	ChannelsRecruiter []Channel `gorm:"foreignKey:RecruiterID"`
	Messages          []Message `gorm:"foreignKey:SenderID"`
}
