package models

import "gorm.io/gorm"

type Lead struct {
	ID                 string `gorm:"primaryKey"`
	Name               string
	ContactInformation string
	LeadSource         string
	InitialContactDate string
	LeadOwner          string
	LeadStage          string
	LeadScore          int32
	Activities         []Activity `gorm:"foreignKey:LeadID"`
}

type Activity struct {
	ID                   string `gorm:"primaryKey"`
	LeadID               string `gorm:"index"`
	ActivityType         string
	DateTime             string
	CommunicationChannel string
	ContentNotes         string
	ParticipantDetails   string
	FollowUpActions      string
}
type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey"`
	GoogleId string
	Name     string
	Email    string
	Phone    string
	Role     string
	Password string
}

var Users []User

var Leads []Lead
