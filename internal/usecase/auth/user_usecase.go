package usecase

import (
	"LoyaltyPointSystem/internal/models"
	auth "LoyaltyPointSystem/internal/repository/auth"
	"LoyaltyPointSystem/utils"
	"errors"
	"log"
)

type AuthUseCase struct {
	userRepo auth.UserRepository
}

func NewAuthUseCase(userRepo auth.UserRepository) *AuthUseCase {
	return &AuthUseCase{userRepo: userRepo}
}

func (uc *AuthUseCase) EmailLogin(email, password string) (*models.User, string, string, error) {
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		return &models.User{}, "", "", errors.New("database error")
	}

	if user.ID == 0 {
		return &models.User{}, "", "", errors.New("invalid email or password")
	}

	// Compare password
	err = utils.CheckPassword(user.PasswordHash, password)
	if err != nil {
		log.Printf("password missmatch issue")
		return &models.User{}, "", "", errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(user.ID)
	if err != nil {
		return &models.User{}, "", "", err
	}
	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return &models.User{}, "", "", err
	}

	err = uc.userRepo.StoreRefreshToken(user.ID, refreshToken)
	if err != nil {
		return &models.User{}, "", "", errors.New("failed to store refresh token")
	}

	return user, accessToken, refreshToken, nil
}

func (uc *AuthUseCase) Refresh(refreshToken string) (uint, string, error) {
	userID, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		return 0, "", err
	}

	storedToken, err := uc.userRepo.GetRefreshToken(userID)
	if err != nil || storedToken != refreshToken {
		return 0, "", errors.New("invalid refresh token")
	}

	newAccessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		return 0, "", err
	}

	return userID, newAccessToken, nil
}
