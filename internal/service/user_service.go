package service

import (
	"errors"

	"github.com/Shubham23jha/go-gin-postgres-clean/internal/models"
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/repository"
	"github.com/Shubham23jha/go-gin-postgres-clean/pkg/utils"
)

type userService struct {
	repo        repository.UserRepository
	sessionRepo repository.SessionRepository
}

func NewUserService(r repository.UserRepository, sr repository.SessionRepository) UserService {
	return &userService{
		repo:        r,
		sessionRepo: sr,
	}
}

func (s *userService) Register(req models.RegisterRequest) error {
	// Optional: check duplicate email
	_, err := s.repo.FindByEmail(req.Email)
	if err == nil {
		return errors.New("email already exists")
	}

	_, err = s.repo.FindByNumber(req.PhoneNumber)
	if err == nil {
		return errors.New("phone number already exists")
	}

	user := models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password, // hashed by hook
		AccountName: req.AccountName,
	}

	return s.repo.Create(&user)
}

func (s *userService) LoginWithPhone(phoneNumber string, password string, deviceID string, deviceName string, browser string, ip string) (string, string, error) {

	user, err := s.repo.GetByPhone(phoneNumber)
	if err != nil {
		return "", "",errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return "","", errors.New("invalid credentials")
	}

	count, _ := s.sessionRepo.CountActive(user.ID)

	if count >= 3 {
		return "", "", errors.New("device limit reached")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "","", err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *userService) LoginWithEmail(email string, password string, deviceID string, deviceName string, browser string, ip string) (string,string, error) {

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "","", errors.New("invalid credentials")
	}

	if !user.CheckPassword(password) {
		return "","", errors.New("invalid credentials")
	}

	count, _ := s.sessionRepo.CountActive(user.ID)

	if count >= 3 {
		return "", "", errors.New("device limit reached")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email)
	if err != nil {
		return "","", err
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Email)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}
