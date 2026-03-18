package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shubham23jha/digital-post-office/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CampaignRepository interface {
	CreateWithOutbox(campaign *models.Campaign, recipients []string) error
	FindByID(id uint) (*models.Campaign, error)
	FetchPendingOutbox(limit int) ([]models.Outbox, error)
	UpdateOutboxStatus(id uint, status models.OutboxStatus) error
	UpdateOutboxFailure(id uint, errMsg string, retryCount int) error
	ReclaimStalledOutbox(timeoutMinutes int) (int64, error)
	CreateEmailLog(log *models.EmailLog) error
	IsMessageProcessed(messageID string) (bool, error)
	FindAll() ([]models.Campaign, error)
}

type campaignRepository struct {
	db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) CampaignRepository {
	return &campaignRepository{db: db}
}

func (r *campaignRepository) CreateWithOutbox(campaign *models.Campaign, recipients []string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Create Campaign
		if err := tx.Create(campaign).Error; err != nil {
			return err
		}

		// 2. Prepare Outbox items
		for _, email := range recipients {
			payload, _ := json.Marshal(map[string]string{
				"campaign_id": fmt.Sprintf("%d", campaign.ID),
				"recipient":   email,
				"subject":     campaign.Subject,
				"body":        campaign.Body,
			})

			outbox := models.Outbox{
				CampaignID: campaign.ID,
				Recipient:  email,
				Payload:    string(payload),
				Status:     models.OutboxStatusPending,
			}

			if err := tx.Create(&outbox).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *campaignRepository) FindByID(id uint) (*models.Campaign, error) {
	var campaign models.Campaign
	err := r.db.Preload("OutboxItems").First(&campaign, id).Error
	return &campaign, err
}

func (r *campaignRepository) FetchPendingOutbox(limit int) ([]models.Outbox, error) {
	var items []models.Outbox
	err := r.db.Clauses(clause.Locking{
		Strength: "UPDATE",
		Options:  "SKIP LOCKED",
	}).
		Where("status = ?", models.OutboxStatusPending).
		Limit(limit).
		Find(&items).Error
	return items, err
}

func (r *campaignRepository) UpdateOutboxStatus(id uint, status models.OutboxStatus) error {
	return r.db.Model(&models.Outbox{}).Where("id = ?", id).Update("status", status).Error
}

func (r *campaignRepository) UpdateOutboxFailure(id uint, errMsg string, retryCount int) error {
	return r.db.Model(&models.Outbox{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      models.OutboxStatusFailed,
		"last_error":  errMsg,
		"retry_count": retryCount,
	}).Error
}

func (r *campaignRepository) ReclaimStalledOutbox(timeoutMinutes int) (int64, error) {
	// Move PICKED_UP items back to PENDING if they haven't been updated for X minutes
	result := r.db.Model(&models.Outbox{}).
		Where("status = ? AND updated_at < ?", models.OutboxStatusPickedUp, time.Now().Add(-time.Duration(timeoutMinutes)*time.Minute)).
		Update("status", models.OutboxStatusPending)
	return result.RowsAffected, result.Error
}

func (r *campaignRepository) CreateEmailLog(emailLog *models.EmailLog) error {
	return r.db.Create(emailLog).Error
}

func (r *campaignRepository) IsMessageProcessed(messageID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.EmailLog{}).Where("message_id = ? AND status = ?", messageID, "SUCCESS").Count(&count).Error
	return count > 0, err
}

func (r *campaignRepository) FindAll() ([]models.Campaign, error) {
	var campaigns []models.Campaign
	err := r.db.Order("created_at desc").Find(&campaigns).Error
	return campaigns, err
}
