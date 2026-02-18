package models

import "time"

type UserSession struct {
	ID uint `gorm:"column:ID;primaryKey;autoIncrement"`

	UserID uint `gorm:"column:userID;not null;index:idx_user_device,priority:1"`

	RefreshToken string `gorm:"column:refreshToken;uniqueIndex;not null"`

	DeviceID   string `gorm:"column:deviceID;index:idx_user_device,priority:2"`
	DeviceName string `gorm:"column:deviceName"`
	Browser    string `gorm:"column:browser"`
	IPAddress  string `gorm:"column:ipAddress"`

	IsActive bool `gorm:"column:isActive;default:true"`

	ExpiresAt time.Time  `gorm:"column:expiresAt"`
	CreatedAt time.Time  `gorm:"column:createdAt;autoCreateTime"`

	// Relation
	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (UserSession) TableName() string {
	return "userSessions"
}

