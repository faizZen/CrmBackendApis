package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

// This is the join table that provides many to many relationship between Campaign and User
type CampaignUser struct {
	CampaignID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}
