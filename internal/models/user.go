package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID uint `gorm:"column:ID;primaryKey;autoIncrement"`

	FirstName string `gorm:"column:firstName"`
	LastName  string `gorm:"column:lastName"`

	Email       string `gorm:"column:email;uniqueIndex:users_email_phone_unique"`
	PhoneNumber string `gorm:"column:phoneNumber;uniqueIndex:users_email_phone_unique"`

	Password string `gorm:"column:password;not null"`

	AccountName string `gorm:"column:accountName"`

	IsVerified bool   `gorm:"column:isVerified;default:false"`
	Role       string `gorm:"column:role;default:user"`

	CreatedAt time.Time `gorm:"column:createdAt;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updatedAt;autoUpdateTime"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	return u.hashPassword()
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("Password") {
		return u.hashPassword()
	}
	return nil
}

func (u *User) hashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword(
		[]byte(u.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	return bcrypt.CompareHashAndPassword(
		[]byte(u.Password),
		[]byte(password),
	) == nil
}

type RegisterRequest struct {
	FirstName   string `json:"firstName" binding:"required"`
	LastName    string `json:"lastName" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phoneNumber" binding:"required"`
	Password    string `json:"password" binding:"required,min=6"`
	AccountName string `json:"accountName"`
}

type LoginRequest struct {
	Email    string `json:"email" `
	PhoneNumber string `json:"phoneNumber"`
	Password string `json:"password" binding:"required"`
	DeviceID    string `json:"deviceID" binding:"required"`
	}