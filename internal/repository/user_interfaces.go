package repository

import "github.com/Shubham23jha/go-gin-postgres-clean/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByNumber(number string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByPhone(phoneNumber string) (*models.User, error)
	MarkVerified(userID uint) error
	
}