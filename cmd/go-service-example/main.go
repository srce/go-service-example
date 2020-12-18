package main

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dzyanis/go-service-example/internal/api"
	"github.com/dzyanis/go-service-example/internal/config"
	"github.com/dzyanis/go-service-example/internal/migrations"
	"github.com/dzyanis/go-service-example/internal/transactions"
	"github.com/dzyanis/go-service-example/internal/users"
	"github.com/dzyanis/go-service-example/internal/wallets"
	"github.com/dzyanis/go-service-example/pkg/boot"
	"github.com/dzyanis/go-service-example/pkg/controllers"
	"github.com/dzyanis/go-service-example/pkg/database"
	"github.com/dzyanis/go-service-example/pkg/logger"
)

func main() {
	ctx := context.Background()

	log := logger.NewLogger(logrus.DebugLevel, os.Stdout)

	runner := boot.NewRunner(ctx, log)

	cnf, err := config.LoadEnv(ctx, ".env")
	if err != nil {
		log.Fatal(err)
	}

	db := database.NewDatabase(cnf.Postgres)
	dbBoot := database.NewBoot(db)
	if err := runner.Try(&dbBoot, 10, time.Second); err != nil {
		log.Fatal(err)
	}

	schema := migrations.NewSchema(db.Write(),
		database.DefaultSchemaName, database.DefaultMigrationTableName)
	migratBoot := migrations.NewBoot(log, schema)
	if err := runner.Once(&migratBoot); err != nil {
		log.Fatal(err)
	}

	jsonHelper := controllers.JSONHelper{}

	useresRepository := users.NewRepository(db)
	usersController := users.NewController(log,
		users.NewService(useresRepository), jsonHelper)

	walletsRepository := wallets.NewRepository(db)
	walletsController := wallets.NewController(log,
		wallets.NewService(walletsRepository), jsonHelper)

	transactionsRepository := transactions.NewRepository(db)
	transactionsService := transactions.NewService(log,
		transactionsRepository, useresRepository, walletsRepository)
	transactionsController := transactions.NewController(log,
		transactionsService, jsonHelper)

	httpServer := api.NewServer(cnf.API,
		usersController, walletsController, transactionsController)
	apiBoot := api.NewBoot(log, httpServer)
	if err := runner.Try(&apiBoot, 3, time.Second); err != nil {
		log.Fatal(err)
	}

	runner.WaitForTermination()
}
