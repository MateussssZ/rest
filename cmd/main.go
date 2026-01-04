package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"rest/config"
	"rest/internal/app"
	"rest/internal/pkg/appLogger"
	"sync"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // Контекст на прерывание
	defer stop()

	conf, err := config.NewConfig() // Читаем конфиг из .env
	if err != nil {
		slog.Error(fmt.Sprintf("configurations initialization failed - error: %v", err))
		return
	}

	logger := appLogger.NewAppLogger(slog.Level(conf.LoggerLevel)) // логгер в консоль

	application, err := app.NewApp(ctx, app.AppDep{ // Инициализируем наше приложение(все бд, http-сервера и тд)
		Config: conf,
		Logger: logger,
	})
	if err != nil {
		logger.Error(ctx, fmt.Errorf("application initialization failed: %w", err))
		return
	}

	wg := &sync.WaitGroup{}
	err = application.Start(ctx, wg) // Запускаем сервер
	if err != nil {
		logger.Error(ctx, fmt.Errorf("application start failed: %w", err))
		return
	}

	logger.Info(ctx, "Application has been started!")
	<-ctx.Done() // Как только наш контекст отменяется прерыванием - делаем graceful shutdown и выходим
	application.Stop(ctx)
	logger.Info(ctx, "Please wait, services are stopping...Chill around 30 seconds")
	wg.Wait()
	logger.Info(ctx, "Application is stopped correctly. The force will be with you")
}
