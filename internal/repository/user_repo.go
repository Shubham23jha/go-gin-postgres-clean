package repository

import (
	"github.com/Shubham23jha/go-gin-postgres-clean/internal/models"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByNumber(number string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("phoneNumber = ?", number).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) GetByPhone(phoneNumber string) (*models.User, error){
	var user models.User
	if err := r.db.Where("phoneNumber = ?", phoneNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
