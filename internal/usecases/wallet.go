package usecases

import (
	"context"
	"rest/internal/repo/postgres"
	"rest/internal/usecases/converts"
	ucmodels "rest/internal/usecases/models"

	"github.com/asaskevich/govalidator"
)

type WalletUsecaseI interface {
	GetWallet(ctx context.Context, walletUUID string) (*ucmodels.Wallet, error)
	WalletOperation(ctx context.Context, operation ucmodels.TransactionRequest) error
}

type WalletUsecaseDep struct {
	Repo postgres.PostgresI
}

type WalletUsecase struct {
	repo postgres.PostgresI
}

func NewWalletUsecase(dep *WalletUsecaseDep) (*WalletUsecase, error) {
	if _, err := govalidator.ValidateStruct(dep); err != nil {
		return nil, err
	}

	return &WalletUsecase{
		repo: dep.Repo,
	}, nil
}

func (u *WalletUsecase) GetWallet(ctx context.Context, walletUUID string) (*ucmodels.Wallet, error) {
	wallet, err := u.repo.GetWalletByID(ctx, walletUUID) // обращаемся к repo-слою
	if err != nil {
		return nil, err
	}

	return converts.WalletRepoToUsecase(wallet), err // http-обработчики не должны иметь никакой связи с repo, поэтому конвертируем repo-структуры в usecase-структуры для изоляции
}

func (u *WalletUsecase) WalletOperation(ctx context.Context, operation ucmodels.TransactionRequest) error {
	err := u.repo.MakeTransaction(ctx, operation.WalletID, operation.OperationType, operation.Amount) // обращаемся к repo-слою
	if err != nil {
		return err
	}

	return nil // конвертировать нечего, поэтому просто сообщаем контроллерам, что ошибок нет
}
