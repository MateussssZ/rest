package errorspkg

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrWalletUUIDIsMissed = errors.New("wallet uuid is missed") // отсутствует uuid кошелька в запросе
	ErrWrongOperationType = errors.New("wrong operation type param") // неправильный тип операции
	ErrWrongAmount        = errors.New("wrong amount param(must be positive)") // отрицательная сумма для операции по кошельку
)

type AppError struct { // наша кастомная ошибка, в которой будет храниться и ошибка, и заранее подобранный к ней статус, также она соответствует интерфейсу error
	Status    int    `json:"-"`
	Err       error  `json:"err"`
	RequestID string `json:"request_id,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s", e.Err.Error())
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func NewBadRequestError(err error) *AppError {
	return &AppError{
		Err:    err,
		Status: http.StatusBadRequest,
	}
}

func NewWalletNotFoundError(walletID string) *AppError {
	return &AppError{
		Err:    fmt.Errorf("Wallet with ID '%s' not found", walletID),
		Status: http.StatusNotFound,
	}
}

func NewInsufficientFundsError(current, required float64) *AppError {
	return &AppError{
		Err:    fmt.Errorf("Insufficient funds, current: %.2f, required: %.2f", current, required),
		Status: http.StatusUnprocessableEntity,
	}
}

func NewInternalError(err error) *AppError {
	return &AppError{
		Status: http.StatusInternalServerError,
		Err:    err,
	}
}

func NewDatabaseError(err error) *AppError {
	return &AppError{
		Status: http.StatusInternalServerError,
		Err:    fmt.Errorf("database error: %w", err),
	}
}
