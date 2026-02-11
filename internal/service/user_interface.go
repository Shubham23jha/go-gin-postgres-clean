package service

import (
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/models"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUsers() ([]models.User, error)
}
