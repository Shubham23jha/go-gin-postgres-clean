package repository

import "github.com/Shubham23jha/digital-post-office/internal/models"

type SessionRepository interface {
	Create(session *models.UserSession) error

	CountActive(userID uint) (int64, error)

	FindByRefreshToken(token string) (*models.UserSession, error)

	DeactivateByToken(token string) error

	DeleteByUserID(userID uint) error
}
