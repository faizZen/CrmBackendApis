package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Lead struct {
	gorm.Model
	ID                 uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FirstName          string    `json:"firstName"`
	LastName           string    `json:"lastName"`
	Email              string    `json:"email"`
	LinkedIn           string    `json:"linkedIn"`
	Country            string    `json:"country"`
	Phone              string    `json:"phone"`
	LeadSource         string    `json:"leadSource"`
	InitialContactDate time.Time `json:"initialContactDate"`
	LeadCreatedBy      uuid.UUID `gorm:"index" json:"leadCreatedBy"`
	Creator            User      `gorm:"foreignKey:LeadCreatedBy;constraint:OnDelete:SET NULL;" json:"creator"`

	// Foreign Key for Assignee
	LeadAssignedTo uuid.UUID `gorm:"index" json:"leadAssignedTo"`
	Assignee       User      `gorm:"foreignKey:LeadAssignedTo;constraint:OnDelete:SET NULL;" json:"assignee"`

	LeadStage      string       `json:"leadStage"`
	LeadNotes      string       `json:"leadNotes"`
	LeadPriority   string       `json:"leadPriority"`
	OrganizationID uuid.UUID    `gorm:"index" json:"organizationId"`
	Organization   Organization `gorm:"foreignKey:OrganizationID" json:"organization"`
	CampaignID     uuid.UUID    `gorm:"index" json:"campaignId"`
	Campaign       Campaign     `gorm:"foreignKey:CampaignID" json:"campaign"`
	Activities     []Activity   `gorm:"foreignKey:LeadID" json:"activities"`
}
type LeadStageHistory struct {
	gorm.Model
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	LeadID    uint      `gorm:"index"`
	OldStage  string
	NewStage  string
	ChangedAt time.Time
}

type Activity struct {
	// ActivityID           string `gorm:"primaryKey" json:"activityId"`
	gorm.Model
	ID                   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	LeadID               uuid.UUID `gorm:"index" json:"leadId"`
	ActivityType         string    `json:"activityType"`
	DateTime             time.Time `json:"dateTime"`
	CommunicationChannel string    `json:"communicationChannel"`
	ContentNotes         string    `json:"contentNotes"`
	ParticipantDetails   string    `json:"participantDetails"`
	FollowUpActions      string    `json:"followUpActions"`
}

type Campaign struct {
	gorm.Model
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CampaignName     string    `json:"campaignName"`
	CampaignCountry  string    `json:"campaignCountry"`
	CampaignRegion   string    `json:"campaignRegion"`
	IndustryTargeted string    `json:"industryTargeted"`
	Leads            []Lead    `gorm:"foreignKey:CampaignID" json:"leads"`
	Users            []User    `gorm:"many2many:campaign_users;joinForeignKey:CampaignID;joinReferences:UserID;constraint:OnDelete:CASCADE;" json:"users"`
}
type CampaignUser struct {
	CampaignID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}
type User struct {
	gorm.Model
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	GoogleId  string     `json:"googleId"`
	Name      string     `json:"name"`
	Email     string     `gorm:"unique" json:"email"`
	Phone     string     `json:"phone"`
	Role      string     `json:"role"`
	Password  string     `json:"password"`
	Campaigns []Campaign `gorm:"many2many:campaign_users;joinForeignKey:UserID;joinReferences:CampaignID;constraint:OnDelete:CASCADE;" json:"campaigns"`
}

type CaseStudy struct {
	gorm.Model
	ID              uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ProjectName     string
	ClientName      string
	TechStack       string
	ProjectDuration string
	KeyOutcomes     string
	IndustryTarget  string
	Tags            string
	Document        string
}

type Organization struct {
	gorm.Model
	ID                  uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	OrganizationName    string    `json:"organizationName"`
	OrganizationEmail   string    `json:"organizationEmail"`
	OrganizationWebsite string    `json:"organizationWebsite"`
	City                string    `json:"city"`
	Country             string    `json:"country"`
	NoOfEmployees       string    `json:"noOfEmployees"`
	AnnualRevenue       string    `json:"annualRevenue"`
	Leads               []Lead    `gorm:"foreignKey:OrganizationID" json:"leads"`
}

type Deals struct {
	gorm.Model
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DealName      string    `json:"dealName"`
	LeadID        uuid.UUID `gorm:"index" json:"leadId"`
	DealStartDate time.Time `json:"dealStartDate"`
	DealEndDate   time.Time `json:"dealEndDate"`

	ProjectRequirements string `json:"projectRequirements"`
	DealAmount          string `json:"dealAmount"`
	DealStatus          string `json:"dealStatus"`
}

type ResourceType string

const (
	ResourceTypeConsultant ResourceType = "CONSULTANT"
	ResourceTypeFreelancer ResourceType = "FREELANCER"
	ResourceTypeContractor ResourceType = "CONTRACTOR"
	ResourceTypeEmployee   ResourceType = "EMPLOYEE"
)

type ResourceStatus string

const (
	ResourceStatusActive   ResourceStatus = "ACTIVE"
	ResourceStatusInactive ResourceStatus = "INACTIVE"
	ResourceStatusOnBench  ResourceStatus = "ON_BENCH"
)

type VendorStatus string

const (
	VendorStatusActive    VendorStatus = "ACTIVE"
	VendorStatusInactive  VendorStatus = "INACTIVE"
	VendorStatusPreferred VendorStatus = "PREFERRED"
)

type PaymentTerms string

const (
	PaymentTermsNet30 PaymentTerms = "NET_30"
	PaymentTermsNet60 PaymentTerms = "NET_60"
	PaymentTermsNet90 PaymentTerms = "NET_90"
)

// type BaseModel struct {
// 	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
// 	CreatedAt time.Time      `gorm:"not null;default:current_timestamp" json:"createdAt"`
// 	UpdatedAt time.Time      `gorm:"not null;default:current_timestamp" json:"updatedAt"`
// 	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
// }

type ResourceProfile struct {
	gorm.Model
	ID                 uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Type               ResourceType    `gorm:"type:resource_type;not null" json:"type"`
	FirstName          string          `gorm:"type:varchar(50);not null" json:"firstName" validate:"min=2,max=50"`
	LastName           string          `gorm:"type:varchar(50);not null" json:"lastName" validate:"min=2,max=50"`
	TotalExperience    float64         `gorm:"not null" json:"totalExperience" validate:"min=0"`
	ContactInformation json.RawMessage `gorm:"type:jsonb;not null" json:"contactInformation"`
	GoogleDriveLink    *string         `gorm:"type:varchar(255)" json:"googleDriveLink,omitempty"`
	Status             ResourceStatus  `gorm:"type:resource_status;not null" json:"status"`
	VendorID           uuid.UUID       `gorm:"type:uuid;index" json:"VendorID,omitempty"`

	// Relationships
	Vendor       *Vendor       `gorm:"foreignKey:VendorID" json:"vendor,omitempty"`
	Skills       []Skill       `gorm:"many2many:resource_skills;" json:"skills"`
	PastProjects []PastProject `gorm:"foreignKey:ResourceProfileID" json:"pastProjects"`
}

type Vendor struct {
	// BaseModel
	gorm.Model
	ID              uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CompanyName     string       `gorm:"type:varchar(100);not null;uniqueIndex" json:"companyName" validate:"min=2,max=100"`
	Status          VendorStatus `gorm:"type:vendor_status;not null" json:"status"`
	PaymentTerms    PaymentTerms `gorm:"type:payment_terms;not null" json:"paymentTerms"`
	Address         string       `gorm:"type:text;not null" json:"address" validate:"max=500"`
	GstOrVatDetails *string      `gorm:"type:varchar(50)" json:"gstOrVatDetails,omitempty" validate:"max=50"`
	Notes           *string      `gorm:"type:text" json:"notes,omitempty" validate:"max=1000"`

	// Relationships
	ContactList        []Contact           `gorm:"foreignKey:VendorID" json:"contactList"`
	Skills             []Skill             `gorm:"many2many:vendor_skills;" json:"skills"`
	PerformanceRatings []PerformanceRating `gorm:"foreignKey:VendorID" json:"performanceRatings"`
	Resources          []ResourceProfile   `gorm:"foreignKey:VendorID" json:"resources"`
}

// --- Supporting Models (for relationships, if needed) ---

// Example supporting model (you might need others depending on your data)
type Skill struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
}

type PastProject struct {
	gorm.Model
	ID                uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ResourceProfileID uuid.UUID `gorm:"type:uuid;index" json:"resourceProfileId"`
	ProjectName       string    `gorm:"type:varchar(100);not null" json:"projectName"`
	Description       string    `gorm:"type:text" json:"description"`
	// Add other project details as needed
}
type Contact struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	VendorID    uuid.UUID `gorm:"type:uuid;index" json:"VendorID"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Email       string    `gorm:"type:varchar(100)" json:"email"`
	PhoneNumber string    `gorm:"type:varchar(20)" json:"phoneNumber"`
	// Add other contact details as needed
}

type PerformanceRating struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	VendorID uuid.UUID `gorm:"type:uuid;index" json:"VendorID"`
	Rating   int       `gorm:"not null" json:"rating"`
	Review   *string   `gorm:"type:text" json:"review,omitempty"`
	// Add other rating details as needed
}
