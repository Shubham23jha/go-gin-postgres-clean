package service

import (
	"errors"
	"time"

	"github.com/Shubham23jha/go-gin-postgres-clean/internal/models"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/repository"
)

type sessionService struct {
	repo repository.SessionRepository
}

func NewSessionService(r repository.SessionRepository) SessionService {

	return &sessionService{
		repo: r,
	}
}

func (s *sessionService) CreateSession(session *models.UserSession) error {

	session.IsActive = true
	session.ExpiresAt = time.Now().Add(7 * 24 * time.Hour)

	return s.repo.Create(session)
}

func (s *sessionService) CheckDeviceLimit(userID uint, limit int64) error {

	count, err := s.repo.CountActive(userID)

	if err != nil {
		return err
	}

	if count >= limit {
		return errors.New("device limit reached")
	}

	return nil
}

func (s *sessionService) ValidateRefreshToken(token string) (*models.UserSession, error) {

	session, err := s.repo.FindByRefreshToken(token)

	if err != nil {
		return nil, errors.New("session not found")
	}

	if !session.IsActive {
		return nil, errors.New("session revoked")
	}

	if session.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}

	return session, nil
}

func (s *sessionService) LogoutByToken(token string) error {

	return s.repo.DeactivateByToken(token)
}

func (s *sessionService) LogoutAll(userID uint) error {

	return s.repo.DeleteByUserID(userID)
}
