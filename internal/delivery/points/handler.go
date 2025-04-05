package pointshandler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"LoyaltyPointSystem/internal/logger"
	pointsuc "LoyaltyPointSystem/internal/usecase/points"

	"github.com/gin-gonic/gin"
)

type PointsHandler struct {
	pointsUC *pointsuc.PointsUseCase
}

func NewPointsHandler(pointsUC *pointsuc.PointsUseCase) *PointsHandler {
	return &PointsHandler{
		pointsUC: pointsUC,
	}
}

func (h *PointsHandler) GetUserLoyaltyHistory(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		userIDStr = c.Query("user_id")
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	authUserIDValue, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	authUserID, ok := authUserIDValue.(uint)
	if !ok || authUserID != uint(userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Pagination parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Max cap to prevent abuse
	}

	resp, err := h.pointsUC.GetUserLoyaltyHistory(uint(userID), page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.LogAudit(&authUserID, "LOYALTY_HISTORY_VIEW", fmt.Sprintf("Page=%d, Limit=%d", page, limit))

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Loyalty history retrieved",
		"meta": gin.H{
			"current_points": resp.CurrentPoints,
			"page":           resp.Page,
			"limit":          resp.Limit,
			"total":          resp.Total,
		},
		"data": resp.History,
	})
}

func (h *PointsHandler) GetFilteredLoyaltyHistory(c *gin.Context) {
	userIDStr := c.Param("user_id")
	if userIDStr == "" {
		userIDStr = c.Query("user_id")
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing user_id"})
		return
	}

	txType := c.Query("tx_type")
	if txType != "" && txType != "earn" && txType != "redeem" && txType != "expire" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tx_type. Use 'CREDIT' or 'DEBIT'"})
		return
	}

	var startDatePtr, endDatePtr *time.Time
	const layout = "2006-01-02"

	startStr := c.Query("start_date")
	if startStr != "" {
		startDate, err := time.Parse(layout, startStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format. Use YYYY-MM-DD"})
			return
		}
		startDatePtr = &startDate
	}

	endStr := c.Query("end_date")
	if endStr != "" {
		endDate, err := time.Parse(layout, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format. Use YYYY-MM-DD"})
			return
		}
		endDatePtr = &endDate
	}

	if startDatePtr != nil && endDatePtr != nil && endDatePtr.Before(*startDatePtr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date cannot be before start_date"})
		return
	}

	resp, err := h.pointsUC.FilterUserLoyaltyHistory(uint(userID), txType, startDatePtr, endDatePtr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
