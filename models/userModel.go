package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserDemo struct {
	gorm.Model
	ID                  uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	GoogleId            string    `gorm:"unique;null" json:"googleId"` // Nullable for traditional login users
	Name                string    `json:"name"`
	Email               string    `gorm:"unique;not null" json:"email"`
	Phone               string    `gorm:"null" json:"phone"` // Nullable in case Google users don’t have a phone
	Role                string    `json:"role"`
	Password            string    `gorm:"null" json:"password"` // Nullable for Google OAuth users
	GoogleRefreshToken  string    `gorm:"null" json:"-"`        // Used only for Google OAuth users (store securely)
	Provider            string    `json:"provider"`             // Google, Facebook, etc.
	BackendRefreshToken string    `gorm:"null;unique" json:"-"` // JWT refresh token for your app (traditional login users & Google users)
	BackendTokenExpiry  time.Time `json:"backendTokenExpiry"`   // Expiry for JWT refresh token
	// Campaigns           []Campaign `gorm:"many2many:campaign_users;joinForeignKey:UserID;joinReferences:CampaignID;constraint:OnDelete:CASCADE;" json:"campaigns"`
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
type GoogleUser struct {
	Provider          string
	Email             string
	Name              string
	FirstName         string
	LastName          string
	NickName          string
	Description       string
	UserID            string
	AvatarURL         string
	Location          string
	AccessToken       string
	AccessTokenSecret string
	RefreshToken      string
	ExpiresAt         time.Time
	IDToken           string
}
