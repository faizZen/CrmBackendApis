package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskPriority string

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
