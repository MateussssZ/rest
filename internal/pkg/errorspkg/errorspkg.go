package errorspkg

import "errors"

var (
	ErrWalletUUIDIsMissed = errors.New("wallet uuid is missed")
	ErrWrongOperationType = errors.New("wrong operation type param")
	ErrWrongAmount        = errors.New("wrong amount param(must be positive)")
)
