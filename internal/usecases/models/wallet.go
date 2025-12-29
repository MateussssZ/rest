package models

import (
	"time"

	"github.com/google/uuid"
)

type Wallet struct {
	ID        uuid.UUID
	Balance   float64
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionRequest struct {
	WalletID      string
	OperationType string
	Amount        float64
}
