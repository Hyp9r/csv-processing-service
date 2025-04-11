package main

import (
	"fmt"
	"net/http"
	"os"

	"database/sql"

	"github.com/Hyp9r/csv-processing-service/domain/product"
	productInfra "github.com/Hyp9r/csv-processing-service/infra/product"
	productPort "github.com/Hyp9r/csv-processing-service/port/rest/product"
	"github.com/Netflix/go-env"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type AppConfig struct {
	Port             string `env:"PORT"`
	PostgresUsername string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDatabase string `env:"POSTGRES_DATABASE"`
	PostgresHost     string `env:"POSTGRES_HOST"`
}

func main() {
	// setup logger
	logger := zerolog.New(os.Stderr)

	// init router
	router := http.NewServeMux()

	// load env variables
	var cfg AppConfig
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		logger.Fatal().Err(err).Msg("error loading enviroment variables")
	}

	// postgres connection
	dbConnStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", cfg.PostgresUsername, cfg.PostgresPassword, cfg.PostgresDatabase, cfg.PostgresHost)
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to postgres")
	}
	defer db.Close()


	//initialize repositories
	productRepo := productInfra.NewProductRepository(db, &logger)

	//initialize services
	productService := product.NewProductService(productRepo, &logger)

	//initialize controllers
	productPort.RegisterRoutes(productService, &logger, router)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: router,
	}

	logger.Info().Msgf("starting up the backend webserver on port: %s", cfg.Port)

	err = server.ListenAndServe()
	if err != nil {
		logger.Err(err).Msg("error while starting up the server")
		panic(err)
	}
}