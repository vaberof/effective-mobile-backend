package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/vaberof/effective-mobile-backend/cmd/carcatalog/docs"
	httproutes "github.com/vaberof/effective-mobile-backend/internal/app/entrypoint/http"
	"github.com/vaberof/effective-mobile-backend/internal/domain/car"
	"github.com/vaberof/effective-mobile-backend/internal/infra/integration/api/apicar"
	"github.com/vaberof/effective-mobile-backend/internal/infra/storage/postgres/pgcar"
	"github.com/vaberof/effective-mobile-backend/pkg/database/postgres"
	"github.com/vaberof/effective-mobile-backend/pkg/http/httpserver"
	apicarhttpclient "github.com/vaberof/effective-mobile-backend/pkg/integration/api/apicar"
	"github.com/vaberof/effective-mobile-backend/pkg/logging/logs"
	"os"
)

var appConfigPaths = flag.String("config.files", "not-found.yaml", "ListByFilters of application config files separated by comma")
var environmentVariablesPath = flag.String("env.vars.file", "not-found.env", "Path to environment variables file")

//	@title			Car Catalog API
//	@version		1.0
//	@description	API Server for Car Catalog Application

// @host		localhost:8000
// @BasePath	/api/v1
func main() {
	flag.Parse()
	if err := loadEnvironmentVariables(); err != nil {
		panic(err)
	}

	appConfig := mustGetAppConfig(*appConfigPaths)

	fmt.Printf("%+v\n", appConfig)

	logger := logs.New(os.Stdout, nil)

	pgManagedDatabase, err := postgres.New(&appConfig.Postgres)
	if err != nil {
		panic(err)
	}

	carApiHttpClient := apicarhttpclient.NewHttpClient(&appConfig.CarApiHttpClient, logger)

	carApiService := apicar.NewCarApiService(carApiHttpClient)

	carStorage := pgcar.NewPgCarStorage(pgManagedDatabase.PostgresDb)
	carService := car.NewCarService(carStorage, carApiService, logger)

	handler := httproutes.NewHandler(carService, logger)

	appServer := httpserver.New(&appConfig.Server, logger)

	handler.InitRoutes(appServer.Router)

	serverExitChannel := appServer.StartAsync()

	<-serverExitChannel
}

func loadEnvironmentVariables() error {
	return godotenv.Load(*environmentVariablesPath)
}
