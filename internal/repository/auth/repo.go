package authrepo

import (
	"LoyaltyPointSystem/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	UpdateUser(user *models.User) error
	StoreRefreshToken(userID uint, refreshToken string) error
	GetRefreshToken(userID uint) (string, error)
	GetUserByID(id uint) (*models.User, error)

	BeginTx() *gorm.DB
	UpdateUserTx(tx *gorm.DB, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) StoreRefreshToken(userID uint, refreshToken string) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", refreshToken).Error
}

func (r *userRepository) GetRefreshToken(userID uint) (string, error) {
	var user models.User
	err := r.db.Select("refresh_token").First(&user, userID).Error
	if err != nil {
		return "", err
	}
	return user.RefreshToken, nil
}

func (r *userRepository) DeleteRefreshToken(userID uint) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("refresh_token", "").Error
}

func (r *userRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *userRepository) UpdateUserTx(tx *gorm.DB, user *models.User) error {
	return tx.Save(user).Error
}

func (r *userRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
