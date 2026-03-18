package repository

import (
	"github.com/Shubham23jha/digital-post-office/internal/models"
	"gorm.io/gorm"
)

type sessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(
	session *models.UserSession,
) error {

	return r.db.
		Create(session).Error
}

func (r *sessionRepository) CountActive(
	userID uint,
) (int64, error) {

	var count int64

	err := r.db.
		Model(&models.UserSession{}).
		Where(
			"user_id = ? AND is_active = true",
			userID,
		).
		Count(&count).Error

	return count, err
}

func (r *sessionRepository) FindByRefreshToken(token string) (*models.UserSession, error) {

	var session models.UserSession

	err := r.db.
		Where("refresh_token = ?", token).
		First(&session).Error

	return &session, err
}

func (r *sessionRepository) DeactivateByToken(
	token string,
) error {

	return r.db.
		Model(&models.UserSession{}).
		Where("refresh_token = ?", token).
		Update("is_active", false).Error
}

func (r *sessionRepository) DeleteByUserID(
	userID uint,
) error {

	return r.db.
		Where("user_id = ?", userID).
		Delete(&models.UserSession{}).Error
}
