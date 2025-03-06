package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
