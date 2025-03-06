package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
