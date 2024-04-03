package main

import (
	"errors"
	"github.com/vaberof/effective-mobile-backend/pkg/config"
	"github.com/vaberof/effective-mobile-backend/pkg/database/postgres"
	"github.com/vaberof/effective-mobile-backend/pkg/http/httpserver"
	"github.com/vaberof/effective-mobile-backend/pkg/integration/api/apicar"
	"os"
)

type AppConfig struct {
	Server           httpserver.ServerConfig
	Postgres         postgres.Config
	CarApiHttpClient apicar.HttpClientConfig
}

func mustGetAppConfig(sources ...string) AppConfig {
	config, err := tryGetAppConfig(sources...)
	if err != nil {
		panic(err)
	}

	if config == nil {
		panic(errors.New("config cannot be nil"))
	}

	return *config
}

func tryGetAppConfig(sources ...string) (*AppConfig, error) {
	if len(sources) == 0 {
		return nil, errors.New("at least 1 source must be set for app config")
	}

	provider := config.MergeConfigs(sources)

	var serverConfig httpserver.ServerConfig
	err := config.ParseConfig(provider, "app.http.server", &serverConfig)
	if err != nil {
		return nil, err
	}

	var postgresConfig postgres.Config
	err = config.ParseConfig(provider, "app.postgres", &postgresConfig)
	if err != nil {
		return nil, err
	}
	postgresConfig.User = os.Getenv("POSTGRES_USER")
	postgresConfig.Password = os.Getenv("POSTGRES_PASSWORD")

	var carApiHttpClientConfig apicar.HttpClientConfig
	err = config.ParseConfig(provider, "app.api.car.http.client", &carApiHttpClientConfig)
	if err != nil {
		return nil, err
	}

	appConfig := AppConfig{
		Server:           serverConfig,
		Postgres:         postgresConfig,
		CarApiHttpClient: carApiHttpClientConfig,
	}

	return &appConfig, nil
}
