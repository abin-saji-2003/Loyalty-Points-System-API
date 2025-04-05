package models

import "time"

type Category struct {
	ID         uint    `gorm:"column:id"`
	Category   string  `json:"category"`
	Multiplier float64 `json:"multiplier"`
}

type LoyaltyHistoryResponse struct {
	CurrentPoints int                  `json:"current_points"`
	History       []LoyaltyPointLedger `json:"history"`
	Page          int                  `json:"page"`
	Limit         int                  `json:"limit"`
	Total         int64                `json:"total"`
}

type FilteredLedgerResponse struct {
	Filters struct {
		TxType    string     `json:"tx_type,omitempty"`
		StartDate *time.Time `json:"start_date,omitempty"`
		EndDate   *time.Time `json:"end_date,omitempty"`
	} `json:"filters"`
	Results []LoyaltyPointLedger `json:"results"`
}
