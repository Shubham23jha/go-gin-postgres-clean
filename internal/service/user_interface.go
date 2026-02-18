package service

import (
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/models"
)

type UserService interface {
	Register(user models.RegisterRequest) error
	LoginWithEmail(email string, password string, deviceID string, deviceName string, browser string, ip string) (string,string,error)
	LoginWithPhone(phoneNumber string, password string,deviceID string, deviceName string, browser string, ip string) (string,string, error)
}
