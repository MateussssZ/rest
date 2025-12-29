package converts

import (
	hndlrmodels "rest/internal/api/web/models"
	ucmodels "rest/internal/usecases/models"
)

func WalletUsecase2Handler(ucWallet *ucmodels.Wallet) hndlrmodels.WalletResponse {
	if ucWallet == nil {
		return hndlrmodels.WalletResponse{}
	}

	return hndlrmodels.WalletResponse{
		ID:       ucWallet.ID,
		Status:   ucWallet.Status,
		Balance:  ucWallet.Balance,
	}
}

func WalletOperationHandler2Usecase(op hndlrmodels.TransactionRequest) ucmodels.TransactionRequest {
	return ucmodels.TransactionRequest{
		WalletID:      op.WalletID,
		OperationType: op.OperationType,
		Amount:        op.Amount,
	}
}
