package models

import "time"

type CampaignStatus string

const (
	CampaignStatusPending   CampaignStatus = "PENDING"
	CampaignStatusRunning   CampaignStatus = "RUNNING"
	CampaignStatusPaused    CampaignStatus = "PAUSED"
	CampaignStatusCompleted CampaignStatus = "COMPLETED"
)

type Campaign struct {
	ID          uint           `gorm:"primaryKey;autoIncrement"`
	Subject     string         `gorm:"size:255;not null"`
	Body        string         `gorm:"type:text;not null"`
	Status      CampaignStatus `gorm:"type:campaign_status;default:'PENDING'"`
	TotalEmails int            `gorm:"default:0"`
	SentEmails  int            `gorm:"default:0"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`

	// Relation
	OutboxItems []Outbox `gorm:"foreignKey:CampaignID;constraint:OnDelete:CASCADE"`
}
