package main

import (
	"fmt"
	"time"

	"LoyaltyPointSystem/db"
	"LoyaltyPointSystem/internal/cronjobs"
	auth "LoyaltyPointSystem/internal/repository/auth"
	points "LoyaltyPointSystem/internal/repository/points"

	"github.com/robfig/cron/v3"
)

func main() {
	fmt.Println("Starting Cron Job Service")

	database := db.ConnectDB()

	pointRepo := points.NewPointRepository(database)
	userRepo := auth.NewUserRepository(database)

	cronService := cronjobs.NewCronService(pointRepo, userRepo)

	c := cron.New()

	_, err := c.AddFunc("0 0 * * *", func() {
		fmt.Println("Running scheduled expiry at", time.Now().Format("2006-01-02 15:04:05"))
		cronService.ExpireOldPoints()
	})
	if err != nil {
		fmt.Println("Failed to add cron job:", err)
		return
	}
	//c.AddFunc("@every 1m", cronService.ExpireOldPoints)

	// Start cron
	c.Start()

	fmt.Println("Cron jobs are running. Waiting for execution...")

	// Keep the service running
	select {}
}
