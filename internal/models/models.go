package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            uint   `gorm:"primaryKey"`
	Email         string `gorm:"unique;not null"`
	PasswordHash  string `gorm:"not null"`
	LoyaltyPoints int    `gorm:"default:0"`
	RefreshToken  string `gorm:"column:refresh_token"`
	CreatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Transaction struct {
	TransactionID   string `gorm:"primaryKey"`
	UserID          uint
	Amount          float64 `gorm:"not null"`
	PaidAmount      float64
	Category        string `gorm:"not null"`
	TransactionDate time.Time
	ProductCode     string
	PointsEarned    int
	CreatedAt       time.Time
}

type RefreshToken struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"type:uuid;not null"`
	Token     string    `gorm:"unique;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}

type LoyaltyPointLedger struct {
	ID            string
	UserID        uint
	TransactionID string
	Points        int
	Type          string // "earn", "redeem", "expired"
	CreatedAt     time.Time
}

type LoyaltyPointRedeemLedger struct {
	ID              string
	UserID          uint
	TransactionID   string
	Status          string
	AvailablePoints int
	ExpiredAt       time.Time
}
