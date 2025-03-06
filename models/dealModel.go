package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Deal struct {
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
