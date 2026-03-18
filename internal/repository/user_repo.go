package repository

import (
	"github.com/Shubham23jha/digital-post-office/internal/models"
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
	if err := r.db.Where("phone_number = ?", number).First(&user).Error; err != nil {
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

func (r *userRepo) GetByPhone(phoneNumber string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) MarkVerified(userID uint) error {

	return r.db.
		Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_verified", true).
		Error
}
