package models

import "gorm.io/gorm"

type Lead struct {
	LeadID             string `gorm:"primaryKey"`
	FirstName          string
	LastName           string
	Email              string
	LinkedIn           string
	Country            string
	Phone              string
	LeadSource         string
	InitialContactDate string
	LeadCreatedBy      string `gorm:"index"`
	Creator            User   `gorm:"foreignKey:LeadCreatedBy;constraint:OnDelete:SET NULL;"`

	// Foreign Key for Assignee
	LeadAssignedTo string `gorm:"index"`
	Assignee       User   `gorm:"foreignKey:LeadAssignedTo;constraint:OnDelete:SET NULL;"`

	LeadStage    string
	LeadNotes    string
	LeadPriority string

	OrganizationID string       `gorm:"index"`
	Organization   Organization `gorm:"foreignKey:OrganizationID"`
	CampaignID     string       `gorm:"index"`
	Campaign       Campaign     `gorm:"foreignKey:CampaignID"`
	Activities     []Activity   `gorm:"foreignKey:LeadID"`
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

type Campaign struct {
	gorm.Model
	CampaignName     string
	CampaignCountry  string
	CampaignRegion   string
	IndustryTargeted string
	// Users            []User `gorm:"many2many:campaign_users"`
	Leads []Lead `gorm:"foreignKey:CampaignID"` // One campaign can have multiple leads

	Users []User `gorm:"many2many:campaign_users;joinForeignKey:CampaignID;joinReferences:UserID;constraint:OnDelete:CASCADE;"`
}

type User struct {
	gorm.Model
	// UserID    string `gorm:"primaryKey"`  // gorm.Model has id field already set
	GoogleId string
	Name     string
	Email    string `gorm:"unique"`
	Phone    string
	Role     string
	Password string
	// Campaigns []Campaign `gorm:"many2many:campaign_users"`
	Campaigns []Campaign `gorm:"many2many:campaign_users;joinForeignKey:UserID;joinReferences:CampaignID;constraint:OnDelete:CASCADE;"`
}

type Organization struct {
	gorm.Model
	// LeadID string `gorm:"index"`
	OrganizationName    string
	OrganizationEmail   string
	OrganizationWebsite string
	City                string
	Country             string
	NoOfEmployees       string
	AnnualRevenue       string
	Leads               []Lead `gorm:"foreignKey:OrganizationID"`
}

// var Users []User

// var Leads []Lead
