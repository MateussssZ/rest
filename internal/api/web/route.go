package rest

import (
	"net/http"
	"rest/internal/api/web/handlers"
	"rest/internal/api/web/middlewares"

	"github.com/asaskevich/govalidator"
	gorillahandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type RouteDeps struct {
	Health        http.Handler
	WalletHandler handlers.WalletI `valid:"required"`
}

func NewRoutes(dep *RouteDeps) (*mux.Router, error) {
	if _, err := govalidator.ValidateStruct(dep); err != nil {
		return nil, err
	}

	r := mux.NewRouter()

	r.Use(gorillahandlers.RecoveryHandler(  // Перехватчик паники
		gorillahandlers.PrintRecoveryStack(true),
	))
	r.Use(middlewares.RequestIDMiddleware) // middleware-обогатитель контекста запроса

	r.Handle("/health", dep.Health) // эндпоинт проверки состояния сервера
	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/wallet", dep.WalletHandler.WalletOperation).Methods(http.MethodPost) // эндпоинт действий с кошельком
	v1.HandleFunc("/wallet/{WALLET_UUID}", dep.WalletHandler.GetWallet).Methods(http.MethodGet) // эндпоинт получения кошелька

	return r, nil
}
