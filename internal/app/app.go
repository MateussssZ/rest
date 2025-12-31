package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"rest/config"
	"rest/internal/pkg/appLogger"

	"github.com/asaskevich/govalidator"
)

type AppDep struct {
	Config *config.Config       `valid:"required"`
	Logger appLogger.IAppLogger `valid:"required"`
}

type App struct {
	srv        *Server
	registries *Registries
}

func NewApp(ctx context.Context, dep AppDep) (*App, error) {
	if _, err := govalidator.ValidateStruct(dep); err != nil {
		return nil, err
	}

	registries, err := NewRepo(ctx, &RepoDep{dep.Config.Postgres})
	if err != nil {
		return nil, err
	}

	usecases, err := NewUsecases(UsecasesDep{
		Repo: registries,
	})
	if err != nil {
		return nil, err
	}

	controllers, err := NewControllers(ControllersDep{
		Usecases: usecases,
	})
	if err != nil {
		return nil, err
	}

	srv, err := NewServer(ServerDep{
		Ctrls:  controllers,
		Config: dep.Config,
		Logger: dep.Logger,
	})
	if err != nil {
		return nil, err
	}

	return &App{
		srv: srv,
	}, nil
}

func (a *App) Start(ctx context.Context, wg *sync.WaitGroup) error {
	errChan := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.srv.Start(ctx)
		if err != nil {
			err = fmt.Errorf("server Started - error: %w", err)
			errChan <- err
		}
	}()

	select {
	case err := <-errChan:
		return err
	case <-time.After(3 * time.Second):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *App) Stop(ctx context.Context) {
	a.srv.Stop(ctx)
	a.registries.Postgres.Stop()
}
