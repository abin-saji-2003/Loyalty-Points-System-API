package transactionuc

import (
	"LoyaltyPointSystem/internal/models"
	authrepo "LoyaltyPointSystem/internal/repository/auth"
	pointrepo "LoyaltyPointSystem/internal/repository/points"
	transactionrepo "LoyaltyPointSystem/internal/repository/transaction"
	"errors"
	"time"

	"github.com/google/uuid"
)

type TransactionUseCase struct {
	transactionRepo transactionrepo.TransactionRepository
	userRepo        authrepo.UserRepository
	pointrepo       pointrepo.PointRepository
}

func NewTransactionUserCase(transactionRepo transactionrepo.TransactionRepository, userRepo authrepo.UserRepository, pointrepo pointrepo.PointRepository) *TransactionUseCase {
	return &TransactionUseCase{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		pointrepo:       pointrepo,
	}
}

func (tc *TransactionUseCase) PlaceTransaction(request models.TransactionRequest) (*models.Transaction, int, error) {
	parsedDate, err := time.Parse("2006-01-02", request.TransactionDate)
	if err != nil {
		return &models.Transaction{}, 0, err
	}

	tx := tc.transactionRepo.BeginTx()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	user, err := tc.userRepo.GetUserByID(request.UserID)
	if err != nil {
		return &models.Transaction{}, 0, err
	}

	discount := 0.0
	pointsEarned := 0
	pointsRedeemed := 0

	if request.UsePoints {
		if user.LoyaltyPoints <= 0 {
			return &models.Transaction{}, 0, errors.New("insufficient points: you do not have enough loyalty points to perform this action")
		}
		maxDiscount := int(request.TransactionAmount * 0.10)
		remaining := maxDiscount

		redeemLedgers, err := tc.pointrepo.GetActiveRedeemLedgersFIFO(tx, request.UserID, parsedDate)
		if err != nil {
			tx.Rollback()
			return &models.Transaction{}, 0, err
		}

		for _, rl := range redeemLedgers {
			if rl.ExpiredAt.Before(parsedDate) || rl.AvailablePoints == 0 || rl.Status == "used" {
				continue
			}

			redeemNow := min(remaining, rl.AvailablePoints)
			rl.AvailablePoints -= redeemNow
			if rl.AvailablePoints == 0 {
				rl.Status = "used"
			}

			err := tc.pointrepo.UpdateRedeemLedgerTx(tx, &rl)
			if err != nil {
				tx.Rollback()
				return &models.Transaction{}, 0, err
			}

			pointsRedeemed += redeemNow
			remaining -= redeemNow
			if remaining == 0 {
				break
			}
		}

		if pointsRedeemed > 0 {
			redeemLedger := &models.LoyaltyPointLedger{
				ID:            uuid.NewString(),
				UserID:        request.UserID,
				TransactionID: request.TransactionID,
				Points:        -pointsRedeemed,
				Type:          "redeem",
				CreatedAt:     parsedDate,
			}
			err := tc.pointrepo.CreateLoyaltyPointLedgerTx(tx, redeemLedger)
			if err != nil {
				tx.Rollback()
				return &models.Transaction{}, 0, err
			}
			user.LoyaltyPoints -= pointsRedeemed
			discount = float64(pointsRedeemed)
		}
	} else {
		category, err := tc.transactionRepo.GetMultiplierByNameTx(tx, request.Category)
		if err != nil {
			tx.Rollback()
			return &models.Transaction{}, 0, err
		}

		finalAmount := request.TransactionAmount
		pointsEarned = int(category.Multiplier * finalAmount)

		expiry := parsedDate.AddDate(1, 0, 0)
		redeemLedger := &models.LoyaltyPointRedeemLedger{
			ID:              uuid.NewString(),
			UserID:          request.UserID,
			TransactionID:   request.TransactionID,
			Status:          "active",
			AvailablePoints: pointsEarned,
			ExpiredAt:       expiry,
		}
		err = tc.pointrepo.CreateRedeemLedgerTx(tx, redeemLedger)
		if err != nil {
			tx.Rollback()
			return &models.Transaction{}, 0, err
		}

		earnLedger := &models.LoyaltyPointLedger{
			ID:            uuid.NewString(),
			UserID:        request.UserID,
			TransactionID: request.TransactionID,
			Points:        pointsEarned,
			Type:          "earn",
			CreatedAt:     parsedDate,
		}
		err = tc.pointrepo.CreateLoyaltyPointLedgerTx(tx, earnLedger)
		if err != nil {
			tx.Rollback()
			return &models.Transaction{}, 0, err
		}
		user.LoyaltyPoints += pointsEarned
	}

	// Save transaction
	transaction := &models.Transaction{
		TransactionID:   request.TransactionID,
		UserID:          request.UserID,
		Amount:          request.TransactionAmount,
		PaidAmount:      request.TransactionAmount - discount,
		Category:        request.Category,
		TransactionDate: parsedDate,
		ProductCode:     request.ProductCode,
		PointsEarned:    pointsEarned,
	}
	_, err = tc.transactionRepo.CreateTransactionTx(tx, transaction)
	if err != nil {
		tx.Rollback()
		return &models.Transaction{}, 0, err
	}

	err = tc.userRepo.UpdateUserTx(tx, user)
	if err != nil {
		tx.Rollback()
		return &models.Transaction{}, 0, err
	}

	if err := tx.Commit().Error; err != nil {
		return &models.Transaction{}, 0, err
	}

	return transaction, user.LoyaltyPoints, nil
}
