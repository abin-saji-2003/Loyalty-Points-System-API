package pointsuc

import (
	"LoyaltyPointSystem/internal/models"
	authrepo "LoyaltyPointSystem/internal/repository/auth"
	pointsrepo "LoyaltyPointSystem/internal/repository/points"
	"time"
)

type PointsUseCase struct {
	pointRepo pointsrepo.PointRepository
	userRepo  authrepo.UserRepository
}

func NewPointsUseCase(pointRepo pointsrepo.PointRepository, userRepo authrepo.UserRepository) *PointsUseCase {
	return &PointsUseCase{
		pointRepo: pointRepo,
		userRepo:  userRepo,
	}
}

func (uc *PointsUseCase) GetUserLoyaltyHistory(userID uint, page, limit int) (*models.LoyaltyHistoryResponse, error) {
	user, err := uc.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	history, total, err := uc.pointRepo.GetFilteredLedgerHistory(userID, page, limit, "", nil, nil)
	if err != nil {
		return nil, err
	}

	response := &models.LoyaltyHistoryResponse{
		CurrentPoints: user.LoyaltyPoints,
		History:       history,
		Page:          page,
		Limit:         limit,
		Total:         total,
	}

	return response, nil
}

func (uc *PointsUseCase) FilterUserLoyaltyHistory(userID uint, txType string, startDate, endDate *time.Time) (*models.FilteredLedgerResponse, error) {
	const maxLimit = 1000

	entries, _, err := uc.pointRepo.GetFilteredLedgerHistory(userID, 1, maxLimit, txType, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var formatted []models.LoyaltyPointLedger
	for _, entry := range entries {
		formatted = append(formatted, models.LoyaltyPointLedger{
			ID:            entry.ID,
			UserID:        entry.UserID,
			TransactionID: entry.TransactionID,
			Points:        entry.Points,
			Type:          entry.Type,
			CreatedAt:     entry.CreatedAt,
		})
	}

	resp := &models.FilteredLedgerResponse{
		Results: formatted,
	}
	resp.Filters.TxType = txType
	resp.Filters.StartDate = startDate
	resp.Filters.EndDate = endDate

	return resp, nil
}
