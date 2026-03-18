// Package service provides the business logic for the digital post office.
package service

import "github.com/Shubham23jha/digital-post-office/internal/models"

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
