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

	deps := postgres.PostgresDeps{
		Conf: &dep.Conf,
	}
	postgresRegistry, err := postgres.New(ctx, deps) // Инициализируем соединение с postgres-базой в другом контейнере
	if err != nil {									// (в структуре Registries может быть сколько угодно реализаций различных БД, которые нам нужны)
		return nil, err
	}

	return &Registries{
		Postgres: postgresRegistry,
	}, nil
}
