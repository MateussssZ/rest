package config

import (
	"github.com/asaskevich/govalidator"
	"github.com/spf13/viper"
)

const _defaultConfigPath = ".env"

type Postgres struct {
	User     string `mapstructure:"POSTGRES_USER" valid:"required"`
	Password string `mapstructure:"POSTGRES_PASSWORD" valid:"required"`
	Host     string `mapstructure:"POSTGRES_HOST" valid:"required"`
	Port     int    `mapstructure:"POSTGRES_PORT" valid:"required"`
	Database string `mapstructure:"POSTGRES_DB" valid:"required"`
}

type Config struct {
	Postgres    Postgres `mapstructure:",squash"`
	Port        int      `mapstructure:"APP_PORT"`
	LoggerLevel int      `mapstructure:"LOGGER_LEVEL"`
}

func NewConfig() (*Config, error) {
	viper.SetConfigFile(_defaultConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var conf Config
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}

	_, err := govalidator.ValidateStruct(conf)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
