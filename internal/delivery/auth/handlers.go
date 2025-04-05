package delivery

import (
	"LoyaltyPointSystem/internal/logger"
	"LoyaltyPointSystem/internal/models"
	usecase "LoyaltyPointSystem/internal/usecase/auth"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUC *usecase.AuthUseCase
}

func NewAuthHandler(authUC *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

func (h *AuthHandler) LoginHandler(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	user, accessToken, refreshToken, err := h.authUC.EmailLogin(req.Email, req.Password)
	if err != nil {
		logger.LogAudit(nil, "LOGIN_FAIL", fmt.Sprintf("Login failed for email: %s - %s", req.Email, err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	logger.LogAudit(&user.ID, "LOGIN", "Successful login")

	c.JSON(http.StatusOK, gin.H{
		"message":       "Login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"user_id":       user.ID,
		"timestamp":     time.Now().Format(time.RFC3339),
	})
}

func (h *AuthHandler) RefreshTokenHandler(c *gin.Context) {
	var req models.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	req.RefreshToken = strings.TrimSpace(req.RefreshToken)
	if req.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Refresh token is required"})
		return
	}

	userID, token, err := h.authUC.Refresh(req.RefreshToken)
	if err != nil {
		logger.LogAudit(nil, "REFRESH_TOKEN_FAIL", fmt.Sprintf("Failed to refresh token: %s", err.Error()))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	logger.LogAudit(&userID, "REFRESH_TOKEN", "Issued new access token")

	c.JSON(http.StatusOK, gin.H{
		"message":      "Access token refreshed successfully",
		"access_token": token,
		"user_id":      userID,
	})
}
