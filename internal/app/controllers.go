package app

import (
	"rest/internal/controllers"

	"github.com/asaskevich/govalidator"
)

type ControllersDep struct {
	Usecases *Usecases `valid:"required"`
}

type Controllers struct {
	Wallet controllers.WalletCtrlI
}

func NewControllers(dep ControllersDep) (*Controllers, error) {
	if _, err := govalidator.ValidateStruct(dep); err != nil {
		return nil, err
	}

	walletDeps := &controllers.WalletCtrlDep{
		WalletUsecase: dep.Usecases.Wallet,
	}

	wallet, err := controllers.NewWalletController(walletDeps)
	if err != nil {
		return nil, err
	}

	return &Controllers{
		Wallet: wallet,
	}, nil
}
