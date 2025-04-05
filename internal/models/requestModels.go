package models

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type TransactionRequest struct {
	TransactionID     string  `json:"transaction_id"`
	UserID            uint    `json:"user_id"`
	TransactionAmount float64 `json:"transaction_amount"`
	Category          string  `json:"category"`
	TransactionDate   string  `json:"transaction_date"`
	ProductCode       string  `json:"product_code"`
	UsePoints         bool    `json:"use_points"`
}
