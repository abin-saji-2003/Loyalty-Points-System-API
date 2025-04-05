package api

import (
	authdelivery "LoyaltyPointSystem/internal/delivery/auth"
	pointshandler "LoyaltyPointSystem/internal/delivery/points"
	transactiondelivery "LoyaltyPointSystem/internal/delivery/transaction"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	authHandler *authdelivery.AuthHandler,
	transactionHandler *transactiondelivery.TransactionHandler,
	pointsHandler *pointshandler.PointsHandler,
	authMiddleware gin.HandlerFunc,
) {
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "server status ok",
		})
	})

	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", authHandler.LoginHandler)
		authGroup.POST("/refresh", authHandler.RefreshTokenHandler)
	}

	transactionGroup := router.Group("/api/transaction")
	transactionGroup.Use(authMiddleware)
	{
		transactionGroup.POST("/", transactionHandler.CreateTransaction)
	}

	pointsGroup := router.Group("/api/points")
	pointsGroup.Use(authMiddleware)
	{
		pointsGroup.GET("/:user_id/history", pointsHandler.GetUserLoyaltyHistory)
		pointsGroup.GET("/:user_id/history/filter", pointsHandler.GetFilteredLoyaltyHistory)
	}
}
