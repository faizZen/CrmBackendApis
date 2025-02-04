package models

import "gorm.io/gorm"

type Lead struct {
	LeadID             string `gorm:"primaryKey"`
	FirstName          string
	LastName           string
	LinkedIn           string
	Country            string
	Phone              string
	LeadSource         string
	InitialContactDate string
	LeadCreatedBy      string `gorm:"index"`
	LeadAssignedTo     string `gorm:"index"`
	LeadStage          string
	LeadNotes          string
	OrganizationID     string `gorm:"index"`

	Activities []Activity `gorm:"foreignKey:LeadID"`
}

type Campaign struct {
	CampaignID       string `gorm:"primaryKey"`
	CampaignName     string
	CampaignCountry  string
	CampaignRegion   string
	IndustryTargeted string
	Users            []User `gorm:"many2many:campaign_users;foreignKey:CampaignID;joinForeignKey:CampaignID;references:UserID;joinReferences:UserID`
}

type Activity struct {
	ActivityID           string `gorm:"primaryKey"`
	LeadID               string `gorm:"index"`
	ActivityType         string
	DateTime             string
	CommunicationChannel string
	ContentNotes         string
	ParticipantDetails   string
	FollowUpActions      string
}

type Organization struct {
	OrganizationID      string `gorm:"primaryKey"`
	OrganizationName    string
	OrganizationEmail   string
	OrganizationWebsite string
	City                string
	Country             string
	NoOfEmployees       string
	AnnualRevenue       string
}

type User struct {
	gorm.Model
	UserID   string `gorm:"primaryKey"`
	GoogleId string
	Name     string
	Email    string
	Phone    string
	Role     string
	Password string
	Campaign string `gorm:"many2many:campaign_users;foreignKey:UserID;joinForeignKey:UserID;references:CampaignID;joinReferences:CampaignID"`
}

var Users []User

var Leads []Lead
