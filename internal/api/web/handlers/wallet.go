package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"rest/internal/api/web/converts"
	hndlrmodels "rest/internal/api/web/models"
	"rest/internal/controllers"
	"rest/internal/pkg/appLogger"
	"rest/internal/pkg/errorspkg"
	httputils "rest/internal/utils/http"

	"github.com/asaskevich/govalidator"
	"github.com/gorilla/mux"
)

type WalletI interface {
	GetWallet(w http.ResponseWriter, r *http.Request)
	WalletOperation(w http.ResponseWriter, r *http.Request)
}

type WalletHandlerDeps struct {
	Logger     appLogger.IAppLogger    `valid:"required"`
	WalletCtrl controllers.WalletCtrlI `valid:"required"`
}

type WalletHandler struct {
	logger     appLogger.IAppLogger
	walletCtrl controllers.WalletCtrlI
}

func NewWalletHandler(dep *WalletHandlerDeps) (*WalletHandler, error) {
	if _, err := govalidator.ValidateStruct(dep); err != nil {
		return nil, err
	}

	return &WalletHandler{
		logger:     dep.Logger,
		walletCtrl: dep.WalletCtrl,
	}, nil
}

func (h *WalletHandler) GetWallet(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	ctx := r.Context()

	vars := mux.Vars(r)
	walletUUID, ok := vars["WALLET_UUID"]
	if !ok {
		logger.Error(ctx, errorspkg.ErrWalletUUIDIsMissed)
		httputils.WriteError(ctx, w, logger, errorspkg.ErrWalletUUIDIsMissed)
		return
	}

	ucWallet, err := h.walletCtrl.GetWallet(ctx, walletUUID)
	if err != nil {
		logger.Error(ctx, err)
		httputils.WriteError(ctx, w, logger, err)
		return
	}

	response := converts.WalletUsecase2Handler(ucWallet)
	httputils.WriteJSON(ctx, w, h.logger, response)
}

func (h *WalletHandler) WalletOperation(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	ctx := r.Context()

	var req hndlrmodels.TransactionRequest
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(ctx, err)
		httputils.WriteError(ctx, w, logger, err)
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		logger.Error(ctx, err)
		httputils.WriteError(ctx, w, logger, err)
		return
	}

	if req.OperationType != "DEPOSIT" && req.OperationType != "WITHDRAW" {
		logger.Error(ctx, errorspkg.ErrWrongOperationType)
		httputils.WriteError(ctx, w, logger, errorspkg.ErrWrongOperationType)
		return
	}

	if req.Amount <= 0 {
		logger.Error(ctx, errorspkg.ErrWrongAmount)
		httputils.WriteError(ctx, w, logger, errorspkg.ErrWrongAmount)
		return
	}

	err = h.walletCtrl.WalletOperation(r.Context(), converts.WalletOperationHandler2Usecase(req))
	if err != nil {
		logger.Error(ctx, err)
		httputils.WriteError(ctx, w, logger, err)
		return
	}

	httputils.WriteJSON(ctx, w, h.logger, "ok")
}
