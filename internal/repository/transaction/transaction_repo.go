package transactionrepo

import (
	"LoyaltyPointSystem/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	BeginTx() *gorm.DB
	CreateTransactionTx(tx *gorm.DB, transaction *models.Transaction) (uint, error)
	GetMultiplierByNameTx(tx *gorm.DB, name string) (models.Category, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (t *transactionRepository) CreateTransactionTx(tx *gorm.DB, transaction *models.Transaction) (uint, error) {
	if err := tx.Create(transaction).Error; err != nil {
		return 0, err
	}
	return transaction.UserID, nil
}

func (t *transactionRepository) GetMultiplierByNameTx(tx *gorm.DB, name string) (models.Category, error) {
	var category models.Category
	err := tx.Where("name = ?", name).First(&category).Error
	if err != nil {
		return models.Category{}, err
	}
	return category, nil
}

func (t *transactionRepository) BeginTx() *gorm.DB {
	return t.db.Begin()
}
