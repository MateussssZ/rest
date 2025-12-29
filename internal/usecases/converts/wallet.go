package converts

import (
	dbmodels "rest/internal/repo/models"
	ucmodels "rest/internal/usecases/models"
)

func WalletRepoToUsecase(dbWallet *dbmodels.Wallet) *ucmodels.Wallet {
	if dbWallet == nil {
		return nil
	}

	return &ucmodels.Wallet{
		ID:        dbWallet.ID,
		Status:    dbWallet.Status,
		Balance:   dbWallet.Balance,
		CreatedAt: dbWallet.CreatedAt,
		UpdatedAt: dbWallet.UpdatedAt,
	}
}
