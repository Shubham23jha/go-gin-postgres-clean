package models

import "time"

type EmailLog struct {
	ID           uint      `gorm:"primaryKey;autoIncrement"`
	CampaignID   uint      `gorm:"not null"`
	MessageID    string    `gorm:"size:255;not null;index"`
	Recipient    string    `gorm:"size:255;not null"`
	Status       string    `gorm:"size:50;not null"` // SUCCESS, FAILED
	ErrorMessage string    `gorm:"type:text"`
	AttemptedAt  time.Time `gorm:"autoCreateTime"`
}
