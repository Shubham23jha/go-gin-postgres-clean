package models

import (
	"time"

	"github.com/google/uuid"
)

type OutboxStatus string

const (
	OutboxStatusPending   OutboxStatus = "PENDING"
	OutboxStatusPickedUp  OutboxStatus = "PICKED_UP"
	OutboxStatusPublished OutboxStatus = "PUBLISHED"
	OutboxStatusFailed    OutboxStatus = "FAILED"
)

type Outbox struct {
	ID         uint         `gorm:"primaryKey;autoIncrement"`
	CampaignID uint         `gorm:"not null;index:idx_outbox_campaign"`
	Recipient  string       `gorm:"size:255;not null"`
	Payload    string       `gorm:"type:jsonb;not null"` // Using string for JSONB mapping
	Status     OutboxStatus `gorm:"type:outbox_status;default:'PENDING';index:idx_outbox_status"`
	RetryCount int          `gorm:"default:0"`
	LastError  string       `gorm:"type:text"`
	MessageID  uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();uniqueIndex"`
	FailedAt   *time.Time
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (Outbox) TableName() string {
	return "outbox"
}
