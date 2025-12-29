package rest

import (
	"net/http"
	"rest/internal/api/web/handlers"

	"github.com/asaskevich/govalidator"
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

	r.Handle("/health", dep.Health)
	v1 := r.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/wallet", dep.WalletHandler.WalletOperation).Methods(http.MethodPost)
	v1.HandleFunc("/wallets/{WALLET_UUID}", dep.WalletHandler.GetWallet).Methods(http.MethodGet)

	return r, nil
}
