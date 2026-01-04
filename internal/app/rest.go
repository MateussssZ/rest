package app

import (
	"context"
	"fmt"
	"net/http"
	"rest/config"
	rest "rest/internal/api/web"
	"rest/internal/api/web/handlers"
	"rest/internal/pkg/appLogger"
	"time"

	"github.com/asaskevich/govalidator"
)

type ServerDep struct {
	Ctrls  *Controllers         `valid:"required"`
	Config *config.Config       `valid:"required"`
	Logger appLogger.IAppLogger `valid:"required"`
}

type Server struct {
	srv *http.Server
}

type webHandlers struct {
	Health        http.Handler
	walletHandler handlers.WalletI
}

type InitHandlersDep struct {
	Logger appLogger.IAppLogger `valid:"required"`
	Ctrls  *Controllers         `valid:"required"`
}

func NewServer(dep ServerDep) (*Server, error) {
	if _, err := govalidator.ValidateStruct(dep); err != nil {
		return nil, err
	}

	handlersDep := &InitHandlersDep{
		Logger: dep.Logger,
		Ctrls:  dep.Ctrls,
	}
	handlers, err := initHandlers(handlersDep) // инициализируем функции-обработчики поступающих запросов
	if err != nil {
		return nil, err
	}

	routesDep := &rest.RouteDeps{
		Health:        handlers.Health,
		WalletHandler: handlers.walletHandler,
	}
	router, err := rest.NewRoutes(routesDep) // инициализируем роутер для нашего сервера
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", "0.0.0.0", dep.Config.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		srv: server,
	}, nil
}

func (s *Server) Start(ctx context.Context) error { // Запуск http-слушателя у сервера
	stop := make(chan error, 1)

	go func() {
		stop <- s.srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return s.Stop(ctx)
	case err := <-stop:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func initHandlers(dep *InitHandlersDep) (*webHandlers, error) {
	walletDeps := &handlers.WalletHandlerDeps{
		Logger:     dep.Logger,
		WalletCtrl: dep.Ctrls.Wallet,
	}
	wallet, err := handlers.NewWalletHandler(walletDeps)
	if err != nil {
		return nil, err
	}

	return &webHandlers{ // обработчики для кошельков и эндпоинта проверки состояния сервера
		Health:        handlers.InitHandlerHealth(),
		walletHandler: wallet,
	}, nil
}
