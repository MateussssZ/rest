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
	handlerName := "GetWallet"
	ctx := r.Context()

	vars := mux.Vars(r)
	walletUUID, ok := vars["WALLET_UUID"]
	if !ok {
		logger.Error(r.Context(), errorspkg.ErrWalletUUIDIsMissed, "handlerName", handlerName)
		httputils.WriteError(ctx, w, logger, http.StatusInternalServerError, errorspkg.ErrWalletUUIDIsMissed)
		return
	}

	ucWallet, err := h.walletCtrl.GetWallet(ctx, walletUUID)
	if err != nil {
		logger.Error(r.Context(), err, "handlerName", handlerName)
		httputils.WriteError(ctx, w, logger, http.StatusInternalServerError, err)
		return
	}

	response := converts.WalletUsecase2Handler(ucWallet)
	httputils.WriteJSON(r.Context(), w, h.logger, response)
}

func (h *WalletHandler) WalletOperation(w http.ResponseWriter, r *http.Request) {
	logger := h.logger
	handlerName := "WalletOperation"
	ctx := r.Context()

	var req hndlrmodels.TransactionRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error(r.Context(), err, "handlerName", handlerName)
		httputils.WriteError(ctx, w, logger, http.StatusInternalServerError, err)
		return
	}

	if err = json.Unmarshal(body, &req); err != nil {
		logger.Error(r.Context(), err, "handlerName", handlerName)
		httputils.WriteError(ctx, w, logger, http.StatusInternalServerError, err)
		return
	}

	if req.OperationType != "DEPOSIT" && req.OperationType != "WITHDRAW" {
		logger.Error(r.Context(), errorspkg.ErrWrongOperationType, "handlerName", handlerName)
		httputils.WriteError(ctx, w, logger, http.StatusInternalServerError, errorspkg.ErrWrongOperationType)
		return
	}

	if req.Amount <= 0 {
		logger.Error(r.Context(), errorspkg.ErrWrongAmount, "handlerName", handlerName)
		httputils.WriteError(ctx, w, logger, http.StatusInternalServerError, errorspkg.ErrWrongAmount)
		return
	}

	err = h.walletCtrl.WalletOperation(r.Context(), converts.WalletOperationHandler2Usecase(req))
	if err != nil {
		logger.Error(r.Context(), err, "handlerName", handlerName)
		httputils.WriteError(ctx, w, logger, http.StatusInternalServerError, err)
		return
	}

	httputils.WriteJSON(r.Context(), w, h.logger, "ok")
}
