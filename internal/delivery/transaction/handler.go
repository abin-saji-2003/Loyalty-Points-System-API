package transaction

import (
	"LoyaltyPointSystem/internal/logger"
	"LoyaltyPointSystem/internal/models"
	usecase "LoyaltyPointSystem/internal/usecase/transaction"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionUC *usecase.TransactionUseCase
}

func NewTransactionHandler(transactionUC *usecase.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{transactionUC: transactionUC}
}

func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	switch {
	case req.TransactionID == "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction ID is required"})
		return
	case req.UserID == 0:
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required and must be valid"})
		return
	case req.TransactionAmount <= 0:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction amount must be greater than 0"})
		return
	case req.Category == "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category is required"})
		return
	case req.ProductCode == "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product code is required"})
		return
	case req.TransactionDate == "":
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction date is required"})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.TransactionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction date must be in YYYY-MM-DD format"})
		return
	}
	req.TransactionDate = parsedDate.Format("2006-01-02")

	transaction, loyaltyPoints, err := h.transactionUC.PlaceTransaction(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	logger.LogAudit(&req.UserID, "TRANSACTION_PLACED", fmt.Sprintf(
		"TransactionID=%s, Amount=%.2f, Category=%s, PointsRemaining=%d",
		req.TransactionID, req.TransactionAmount, req.Category, loyaltyPoints,
	))

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Transaction recorded and loyalty points updated",
		"data": gin.H{
			"transaction": gin.H{
				"transaction_id":   transaction.TransactionID,
				"user_id":          transaction.UserID,
				"amount":           transaction.Amount,
				"paid_amount":      transaction.PaidAmount,
				"category":         transaction.Category,
				"transaction_date": transaction.TransactionDate.Format(time.RFC3339),
				"product_code":     transaction.ProductCode,
				"points_earned":    transaction.PointsEarned,
				"created_at":       transaction.CreatedAt.Format(time.RFC3339),
			},
			"remaining_loyalty_points": loyaltyPoints,
			"timestamp":                time.Now().Format(time.RFC3339),
		},
	})
}
