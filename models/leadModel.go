package models

import (
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

// Complimentary model in context with lead
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
