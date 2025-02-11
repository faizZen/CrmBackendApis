package models

import "gorm.io/gorm"

type Lead struct {
	LeadID             string         `gorm:"primaryKey" json:"leadId"`
	FirstName          string         `json:"firstName"`
	LastName           string         `json:"lastName"`
	Email              string         `json:"email"`
	LinkedIn           string         `json:"linkedIn"`
	Country            string         `json:"country"`
	Phone              string         `json:"phone"`
	LeadSource         string         `json:"leadSource"`
	InitialContactDate string         `json:"initialContactDate"`
	DeletedAt          gorm.DeletedAt `gorm:"index"` // Soft delete

	LeadCreatedBy string `gorm:"index" json:"leadCreatedBy"`
	Creator       User   `gorm:"foreignKey:LeadCreatedBy;constraint:OnDelete:SET NULL;" json:"creator"`

	// Foreign Key for Assignee
	LeadAssignedTo string `gorm:"index" json:"leadAssignedTo"`
	Assignee       User   `gorm:"foreignKey:LeadAssignedTo;constraint:OnDelete:SET NULL;" json:"assignee"`

	LeadStage      string       `json:"leadStage"`
	LeadNotes      string       `json:"leadNotes"`
	LeadPriority   string       `json:"leadPriority"`
	OrganizationID string       `gorm:"index" json:"organizationId"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	CampaignID     string       `gorm:"index" json:"campaignId"`
	Campaign       Campaign     `gorm:"foreignKey:CampaignID" json:"campaign"`
	Activities     []Activity   `gorm:"foreignKey:LeadID" json:"activities"`
}

type Activity struct {
	ActivityID           string `gorm:"primaryKey" json:"activityId"`
	LeadID               string `gorm:"index" json:"leadId"`
	ActivityType         string `json:"activityType"`
	DateTime             string `json:"dateTime"`
	CommunicationChannel string `json:"communicationChannel"`
	ContentNotes         string `json:"contentNotes"`
	ParticipantDetails   string `json:"participantDetails"`
	FollowUpActions      string `json:"followUpActions"`
}

type Campaign struct {
	gorm.Model
	CampaignName     string `json:"campaignName"`
	CampaignCountry  string `json:"campaignCountry"`
	CampaignRegion   string `json:"campaignRegion"`
	IndustryTargeted string `json:"industryTargeted"`
	Leads            []Lead `gorm:"foreignKey:CampaignID" json:"leads"`
	Users            []User `gorm:"many2many:campaign_users;joinForeignKey:CampaignID;joinReferences:UserID;constraint:OnDelete:CASCADE;" json:"users"`
}

type User struct {
	gorm.Model
	GoogleId  string     `json:"googleId"`
	Name      string     `json:"name"`
	Email     string     `gorm:"unique" json:"email"`
	Phone     string     `json:"phone"`
	Role      string     `json:"role"`
	Password  string     `json:"password"`
	Campaigns []Campaign `gorm:"many2many:campaign_users;joinForeignKey:UserID;joinReferences:CampaignID;constraint:OnDelete:CASCADE;" json:"campaigns"`
}

type Organization struct {
	gorm.Model
	OrganizationName    string `json:"organizationName"`
	OrganizationEmail   string `json:"organizationEmail"`
	OrganizationWebsite string `json:"organizationWebsite"`
	City                string `json:"city"`
	Country             string `json:"country"`
	NoOfEmployees       string `json:"noOfEmployees"`
	AnnualRevenue       string `json:"annualRevenue"`
	Leads               []Lead `gorm:"foreignKey:OrganizationID" json:"leads"`
}

type Deals struct {
	gorm.Model
	DealName            string `json:"dealName"`
	LeadID              string `gorm:"index" json:"leadId"`
	DealStartDate       string `json:"dealStartDate"`
	DealEndDate         string `json:"dealEndDate"`
	ProjectRequirements string `json:"projectRequirements"`
	DealAmount          string `json:"dealAmount"`
	DealStatus          string `json:"dealStatus"`
}
