package main

import (
	"LoyaltyPointSystem/db"
	"LoyaltyPointSystem/internal/api"
	authdelivery "LoyaltyPointSystem/internal/delivery/auth"
	pointshandler "LoyaltyPointSystem/internal/delivery/points"
	transactiondelivery "LoyaltyPointSystem/internal/delivery/transaction"

	"LoyaltyPointSystem/internal/middleware"
	authRepo "LoyaltyPointSystem/internal/repository/auth"
	pointrepo "LoyaltyPointSystem/internal/repository/points"
	transactionRepo "LoyaltyPointSystem/internal/repository/transaction"

	authUC "LoyaltyPointSystem/internal/usecase/auth"
	pointsUC "LoyaltyPointSystem/internal/usecase/points"
	transactionUC "LoyaltyPointSystem/internal/usecase/transaction"

	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the DB
	database := db.ConnectDB()

	// Initialize Repositories
	userRepo := authRepo.NewUserRepository(database)
	txRepo := transactionRepo.NewTransactionRepository(database)
	pointRepo := pointrepo.NewPointRepository(database)

	// Initialize Use Cases
	authUseCase := authUC.NewAuthUseCase(userRepo)
	transactionUseCase := transactionUC.NewTransactionUserCase(txRepo, userRepo, pointRepo)
	pointsUseCase := pointsUC.NewPointsUseCase(pointRepo, userRepo)

	// Initialize Handlers
	authHandler := authdelivery.NewAuthHandler(authUseCase)
	transactionHandler := transactiondelivery.NewTransactionHandler(transactionUseCase)
	pointsHandler := pointshandler.NewPointsHandler(pointsUseCase)

	// Middleware
	authMiddleware := middleware.AuthMiddleware()

	// Initialize Router
	router := gin.Default()

	// Register Routes
	api.RegisterRoutes(router, authHandler, transactionHandler, pointsHandler, authMiddleware)

	port := os.Getenv("PORT")

	fmt.Println("Server running on port:", port)
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
