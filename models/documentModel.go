package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Document struct {
	gorm.Model
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Title         string         `gorm:"size:255;not null" json:"title"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null" json:"userId"`
	FilePath      string         `gorm:"type:varchar(255);not null" json:"filePath"`
	FileSize      string         `gorm:"type:varchar(50);not null" json:"fileSize"`
	FileType      string         `gosrm:"type:varchar(50);not null" json:"fileType"`
	ReferenceID   uuid.UUID      `gorm:"type:uuid;not null" json:"reference"`
	ReferenceType string         `gorm:"type:varchar(50);not null" json:"referenceType"`
	Tags          pq.StringArray `gorm:"type:text[]"`
}