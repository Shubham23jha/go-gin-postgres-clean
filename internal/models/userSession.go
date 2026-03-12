package models

import "time"

type UserSession struct {
	ID uint `gorm:"primaryKey;autoIncrement"`

	UserID uint `gorm:"not null;index:idx_user_device,priority:1"`

	RefreshToken string `gorm:"uniqueIndex;not null"`

	DeviceID   string `gorm:"index:idx_user_device,priority:2"`
	DeviceName string
	Browser    string
	IPAddress  string

	IsActive bool `gorm:"default:true"`

	ExpiresAt time.Time
	CreatedAt time.Time `gorm:"autoCreateTime"`

	// Relation
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
