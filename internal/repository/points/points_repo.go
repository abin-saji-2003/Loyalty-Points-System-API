package pointsrepo

import (
	"LoyaltyPointSystem/internal/models"
	"time"

	"gorm.io/gorm"
)

type PointRepository interface {
	GetFilteredLedgerHistory(userID uint, page, limit int, txType string, startDate, endDate *time.Time) ([]models.LoyaltyPointLedger, int64, error)
	UpdateLedgerTypeToExpired(transactionID string) error
	GetExpiredActiveRedeemLedgers(expireBefore time.Time) ([]models.LoyaltyPointRedeemLedger, error)

	BeginTx() *gorm.DB
	CreateLoyaltyPointLedgerTx(tx *gorm.DB, ledger *models.LoyaltyPointLedger) error
	CreateRedeemLedgerTx(tx *gorm.DB, redeem *models.LoyaltyPointRedeemLedger) error
	GetActiveRedeemLedgersFIFO(tx *gorm.DB, userID uint, currentDate time.Time) ([]models.LoyaltyPointRedeemLedger, error)
	UpdateRedeemLedgerTx(tx *gorm.DB, ledger *models.LoyaltyPointRedeemLedger) error
}

type pointRepository struct {
	db *gorm.DB
}

func NewPointRepository(db *gorm.DB) PointRepository {
	return &pointRepository{db: db}
}

func (r *pointRepository) GetFilteredLedgerHistory(userID uint, page, limit int, txType string, startDate, endDate *time.Time) ([]models.LoyaltyPointLedger, int64, error) {
	var history []models.LoyaltyPointLedger
	var total int64
	offset := (page - 1) * limit

	query := r.db.Model(&models.LoyaltyPointLedger{}).Where("user_id = ?", userID)

	if txType != "" {
		query = query.Where("type = ?", txType)
	}

	// Optional date filters
	if startDate != nil && endDate != nil {
		query = query.Where("created_at BETWEEN ? AND ?", *startDate, *endDate)
	} else if startDate != nil {
		query = query.Where("created_at >= ?", *startDate)
	} else if endDate != nil {
		query = query.Where("created_at <= ?", *endDate)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&history).Error; err != nil {
		return nil, 0, err
	}

	return history, total, nil
}

func (r *pointRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

func (r *pointRepository) CreateLoyaltyPointLedgerTx(tx *gorm.DB, ledger *models.LoyaltyPointLedger) error {
	return tx.Create(ledger).Error
}

func (r *pointRepository) CreateRedeemLedgerTx(tx *gorm.DB, redeem *models.LoyaltyPointRedeemLedger) error {
	return tx.Create(redeem).Error
}

func (r *pointRepository) UpdateRedeemLedgerTx(tx *gorm.DB, ledger *models.LoyaltyPointRedeemLedger) error {
	if tx == nil {
		tx = r.db
	}

	err := tx.Model(&models.LoyaltyPointRedeemLedger{}).
		Where("id = ?", ledger.ID).
		Updates(map[string]interface{}{
			"available_points": ledger.AvailablePoints,
			"status":           ledger.Status,
		}).Error

	return err
}

func (r *pointRepository) GetActiveRedeemLedgersFIFO(tx *gorm.DB, userID uint, currentDate time.Time) ([]models.LoyaltyPointRedeemLedger, error) {
	var ledgers []models.LoyaltyPointRedeemLedger

	err := tx.
		Where("user_id = ? AND status = ? AND expired_at > ? AND available_points > 0", userID, "active", currentDate).
		Order("expired_at ASC").
		Find(&ledgers).Error

	if err != nil {
		return nil, err
	}

	return ledgers, nil
}

func (r *pointRepository) GetExpiredActiveRedeemLedgers(expireBefore time.Time) ([]models.LoyaltyPointRedeemLedger, error) {
	var ledgers []models.LoyaltyPointRedeemLedger

	err := r.db.
		Where("status = ? AND expired_at <= ? AND available_points > 0", "active", expireBefore).
		Find(&ledgers).Error

	if err != nil {
		return nil, err
	}

	return ledgers, nil
}

func (r *pointRepository) UpdateLedgerTypeToExpired(transactionID string) error {
	return r.db.Model(&models.LoyaltyPointLedger{}).
		Where("transaction_id = ? AND type = ?", transactionID, "redeem").
		Update("type", "expired").Error
}
