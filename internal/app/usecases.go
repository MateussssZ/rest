package app

import (
	"rest/internal/usecases"

	"github.com/asaskevich/govalidator"
)

type UsecasesDep struct {
	Repo *Registries `valid:"required"`
}

type Usecases struct {
	Wallet usecases.WalletUsecaseI
}

func NewUsecases(dep UsecasesDep) (*Usecases, error) {
	if _, err := govalidator.ValidateStruct(dep); err != nil {
		return nil, err
	}

	walletDeps := &usecases.WalletUsecaseDep{
		Repo: dep.Repo.Postgres,
	}
	wallet, err := usecases.NewWalletUsecase(walletDeps) // Инициализируем usecase для кошельков(в будущем при появлении других сущностей у нас будут появляться другие usecase`ы)
	if err != nil {
		return nil, err
	}

	return &Usecases{
		Wallet: wallet,
	}, nil
}
