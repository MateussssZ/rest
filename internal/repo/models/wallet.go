package models

import (
	"time"

	"github.com/google/uuid"
)

type WalletStatus string
type OperationType string

const (
	WalletStatusActive WalletStatus = "ACTIVE"
	WalletStatusFrozen WalletStatus = "FROZEN"
	WalletStatusClosed WalletStatus = "CLOSED"
)

const (
	OperationDeposit  OperationType = "DEPOSIT"
	OperationWithdraw OperationType = "WITHDRAW"
)

type Wallet struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Balance   float64   `db:"balance" json:"balance"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type Transaction struct {
	ID              uuid.UUID `db:"id" json:"id"`
	WalletID        uuid.UUID `db:"wallet_id" json:"wallet_id"`
	OperationType   string    `db:"operation_type" json:"operation_type"`
	Amount          float64   `db:"amount" json:"amount"`
	PreviousBalance float64   `db:"previous_balance" json:"previous_balance"`
	NewBalance      float64   `db:"new_balance" json:"new_balance"`
	Status          string    `db:"status" json:"status"`
	Description     *string   `db:"description" json:"description,omitempty"`
	CreatedAt       time.Time `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time `db:"updated_at" json:"updated_at"`
}

type User struct {
	ID        uuid.UUID `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Phone     string    `db:"phone" json:"phone"`
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
