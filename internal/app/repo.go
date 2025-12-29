package app

import (
	"context"
	"rest/config"

	"rest/internal/repo/postgres"

	"github.com/asaskevich/govalidator"
)

type RepoDep struct {
	Conf config.Postgres `valid:"required"`
}

type Registries struct {
	Postgres postgres.PostgresI
}

func NewRepo(ctx context.Context, dep *RepoDep) (*Registries, error) {
	if _, err := govalidator.ValidateStruct(dep); err != nil {
		return nil, err
	}

	// инициализации конекшенов ко всем используемым БД
	deps := postgres.PostgresDeps{
		Conf: &dep.Conf,
	}
	postgresRegistry, err := postgres.New(ctx, deps)
	if err != nil {
		return nil, err
	}

	return &Registries{
		Postgres: postgresRegistry,
	}, nil
}
