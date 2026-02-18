package service

import "github.com/Shubham23jha/go-gin-postgres-clean/internal/models"

type SessionService interface {
	CreateSession(

		session *models.UserSession,
	) error

	CheckDeviceLimit(

		userID uint,
		limit int64,
	) error

	ValidateRefreshToken(

		token string,
	) (*models.UserSession, error)

	LogoutByToken(

		token string,
	) error

	LogoutAll(

		userID uint,
	) error
}
