package controllers

import (
	"context"
	"rest/internal/usecases"
	ucmodels "rest/internal/usecases/models"

	"github.com/asaskevich/govalidator"
)

type WalletCtrlI interface {
	GetWallet(ctx context.Context, walletUUID string) (*ucmodels.Wallet, error)
	WalletOperation(ctx context.Context, op ucmodels.TransactionRequest) error
}

type WalletCtrlDep struct {
	WalletUsecase usecases.WalletUsecaseI
}

type WalletCtrl struct {
	walletUsecase usecases.WalletUsecaseI
}

func NewWalletController(dep *WalletCtrlDep) (*WalletCtrl, error) {
	if _, err := govalidator.ValidateStruct(dep); err != nil {
		return nil, err
	}

	return &WalletCtrl{
		walletUsecase: dep.WalletUsecase,
	}, nil
}

func (u *WalletCtrl) GetWallet(ctx context.Context, walletUUID string) (*ucmodels.Wallet, error) {
	// Здесь в будущем будет логика проверки токенов и прочего, но пока пусто :)

	wallet, err := u.walletUsecase.GetWallet(ctx, walletUUID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (u *WalletCtrl) WalletOperation(ctx context.Context, op ucmodels.TransactionRequest) error {
	// Здесь в будущем будет логика проверки токенов и прочего, но пока пусто :)

	err := u.walletUsecase.WalletOperation(ctx, op)
	if err != nil {
		return err
	}

	return nil
}
