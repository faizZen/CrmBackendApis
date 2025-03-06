package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
	ResourceProfileID uuid.UUID `gorm:"type:uuid;primaryKey" json:"resourceProfileId"`
	SkillID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"skillId"`

	ExperienceYears float64 `gorm:"not null" json:"experienceYears"`

	// Ensure Skill relation is properly mapped
	Skill Skill `gorm:"foreignKey:SkillID;references:ID;constraint:OnDelete:CASCADE;" json:"skill"`

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

type SkillType string

const (
	FRONTEND SkillType = "FRONTEND"
	BACKEND  SkillType = "BACKEND"
	DESIGN   SkillType = "DESIGN"
	// Add other skill types as needed
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"not null"`
	Token     string    `gorm:"unique;not null"`
	CreatedAt time.Time
	ExpiresAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete support
}
