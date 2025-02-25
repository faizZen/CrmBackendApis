package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
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

	LeadStage      LeadStage    `json:"leadStage"`
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
	LeadID    uuid.UUID `gorm:"type:uuid;not null;index" json:"lead_id"`
	Lead      Lead      `gorm:"foreignKey:LeadID;constraint:OnDelete:CASCADE;" json:"lead"`
	OldStage  LeadStage `json:"oldStage"`
	NewStage  LeadStage `json:"newStage"`
	ChangedAt time.Time `json:"changedAt"`
}

type LeadStage string

const (
	LeadStageNew        LeadStage = "NEW"
	LeadStageInProgress LeadStage = "IN_PROGRESS"
	LeadStageFollowUp   LeadStage = "FOLLOW_UP"
	LeadStageClosedWon  LeadStage = "CLOSED_WON"
	LeadStageClosedLost LeadStage = "CLOSED_LOST"
)

type Activity struct {
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
	ID                  uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	DealName            string    `json:"dealName"`
	LeadID              uuid.UUID `gorm:"index" json:"leadId"`
	DealStartDate       time.Time `json:"dealStartDate"`
	DealEndDate         time.Time `json:"dealEndDate"`
	ProjectRequirements string    `json:"projectRequirements"`
	DealAmount          string    `json:"dealAmount"`
	DealStatus          string    `json:"dealStatus"`
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
	VendorID           uuid.UUID       `gorm:"type:uuid;index" json:"vendorId,omitempty"`

	// Instead of a simple many2many, we now use our custom join table.
	ResourceSkills []ResourceSkill `gorm:"foreignKey:ResourceProfileID" json:"resourceSkills"`
	PastProjects   []PastProject   `gorm:"foreignKey:ResourceProfileID" json:"pastProjects"`
}

type ResourceSkill struct {
	// Composite Primary Key: ResourceProfileID and SkillID
	ResourceProfileID uuid.UUID `gorm:"type:uuid;primaryKey" json:"resourceProfileId"`
	SkillID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"skillId"`

	// Extra column for the number of years of experience with this skill.
	ExperienceYears float64 `gorm:"not null" json:"experienceYears"`

	// Association: preload the related Skill record.
	Skill Skill `gorm:"foreignKey:SkillID;references:ID" json:"skill"`

	// Timestamps (optional)
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
	SkillType   SkillType `gorm:"type:varchar(50);not null" json:"skillType"` // New field
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

type TaskPriority string

type SkillType string

const (
	FRONTEND SkillType = "FRONTEND"
	BACKEND  SkillType = "BACKEND"
	DESIGN   SkillType = "DESIGN"
	// Add other skill types as needed
)

const (
	LOW    TaskPriority = "LOW"
	MEDIUM TaskPriority = "MEDIUM"
	HIGH   TaskPriority = "HIGH"
	URGENT TaskPriority = "URGENT"
)

type TaskStatus string

const (
	TODO        TaskStatus = "TODO"
	IN_PROGRESS TaskStatus = "IN_PROGRESS"
	COMPLETED   TaskStatus = "COMPLETED"
	ON_HOLD     TaskStatus = "ON_HOLD"
)

type Task struct {
	gorm.Model
	ID          uuid.UUID    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      uuid.UUID    `gorm:"type:uuid;not null" json:"userId"`
	Title       string       `gorm:"size:255;not null" json:"title"`
	Description string       `gorm:"type:text" json:"description"`
	Status      TaskStatus   `gorm:"type:task_status;not null" json:"status"`
	Priority    TaskPriority `gorm:"type:task_priority;not null" json:"priority"`
	DueDate     *time.Time   `json:"dueDate"` // Nullable time for dueDate
	User        User         `gorm:"foreignKey:UserID;references:ID" json:"user"`
}

type Document struct {
	gorm.Model
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null" json:"userId"`
	FilePath      string         `gorm:"type:varchar(255);not null" json:"filePath"`
	FileSize      string         `gorm:"type:varchar(50);not null" json:"fileSize"`
	FileType      string         `gorm:"type:varchar(50);not null" json:"fileType"`
	ReferenceID   uuid.UUID      `gorm:"type:uuid;not null" json:"reference"`
	ReferenceType string         `gorm:"type:varchar(50);not null" json:"referenceType"`
	Tags          pq.StringArray `gorm:"type:text[]"`
}
