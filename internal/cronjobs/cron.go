package cronjobs

import (
	auth "LoyaltyPointSystem/internal/repository/auth"
	points "LoyaltyPointSystem/internal/repository/points"
	"fmt"
	"time"
)

type CronService struct {
	pointRepo points.PointRepository
	userRepo  auth.UserRepository
}

func NewCronService(pointRepo points.PointRepository, userRepo auth.UserRepository) *CronService {
	return &CronService{
		pointRepo: pointRepo,
		userRepo:  userRepo,
	}
}

func (cs *CronService) ExpireOldPoints() {
	now := time.Now()

	expiredLedgers, err := cs.pointRepo.GetExpiredActiveRedeemLedgers(now)
	if err != nil {
		fmt.Println("Failed to fetch expired redeem ledgers:", err)
		return
	}

	expiredPointsByUser := make(map[uint]int)

	for _, rl := range expiredLedgers {
		if rl.AvailablePoints <= 0 {
			continue
		}

		expiredPoints := rl.AvailablePoints
		rl.Status = "expired"
		rl.AvailablePoints = 0

		if err := cs.pointRepo.UpdateRedeemLedgerTx(nil, &rl); err != nil {
			fmt.Println("Failed to mark redeem ledger as expired:", err)
			continue
		}

		if err := cs.pointRepo.UpdateLedgerTypeToExpired(rl.TransactionID); err != nil {
			fmt.Println("Failed to update ledger entry to expired:", err)
			continue
		}

		expiredPointsByUser[rl.UserID] += expiredPoints
	}

	for userID, expiredPoints := range expiredPointsByUser {
		user, err := cs.userRepo.GetUserByID(userID)
		if err != nil {
			fmt.Printf("User %d not found: %v\n", userID, err)
			continue
		}

		if user.LoyaltyPoints >= expiredPoints {
			user.LoyaltyPoints -= expiredPoints
		} else {
			user.LoyaltyPoints = 0
		}

		if err := cs.userRepo.UpdateUser(user); err != nil {
			fmt.Printf("Failed to update user %d: %v\n", userID, err)
		}
	}

	fmt.Printf("Marked %d redeem ledgers as expired and updated %d users.\n", len(expiredLedgers), len(expiredPointsByUser))
}
