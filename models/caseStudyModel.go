package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
