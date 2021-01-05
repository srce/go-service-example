package main

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dzyanis/go-service-example/internal/api"
	"github.com/dzyanis/go-service-example/internal/config"
	"github.com/dzyanis/go-service-example/internal/migrations"
	transactionsControllers "github.com/dzyanis/go-service-example/internal/transactions/controllers"
	transactionsRepositories "github.com/dzyanis/go-service-example/internal/transactions/repositories"
	transactionsServices "github.com/dzyanis/go-service-example/internal/transactions/services"
	"github.com/dzyanis/go-service-example/internal/transactions/uow"
	usersControllers "github.com/dzyanis/go-service-example/internal/users/controllers"
	usersRepositories "github.com/dzyanis/go-service-example/internal/users/repositories"
	usersServices "github.com/dzyanis/go-service-example/internal/users/services"
	walletsControllers "github.com/dzyanis/go-service-example/internal/wallets/controllers"
	walletsRepositories "github.com/dzyanis/go-service-example/internal/wallets/repositories"
	walletsServices "github.com/dzyanis/go-service-example/internal/wallets/services"
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

	db := database.NewPostgres(cnf.Postgres)
	dbBoot := database.NewBoot(db)
	if err := runner.Try(&dbBoot, 10, time.Second); err != nil {
		log.Fatal(err)
	}

	schema := migrations.NewSchema(db.Connection(),
		database.DefaultSchemaName, database.DefaultMigrationTableName)
	migratBoot := migrations.NewBoot(log, schema)
	if err := runner.Once(&migratBoot); err != nil {
		log.Fatal(err)
	}

	jsonHelper := controllers.JSONHelper{}

	usersRepository := usersRepositories.NewRepository(db)
	usersController := usersControllers.NewController(log,
		usersServices.NewService(usersRepository), jsonHelper)

	walletsRepository := walletsRepositories.NewRepository(db)
	walletsController := walletsControllers.NewController(log,
		walletsServices.NewService(walletsRepository), jsonHelper)

	startUOW := func() (*uow.UOW, error) {
		return uow.NewUOW(db.Connection())
	}

	transactionsRepository := transactionsRepositories.NewRepository(db)
	transactionsService := transactionsServices.NewService(log,
		transactionsRepository, usersRepository, walletsRepository, startUOW)
	transactionsController := transactionsControllers.NewController(log,
		transactionsService, jsonHelper)

	httpServer := api.NewServer(cnf.API,
		usersController, walletsController, transactionsController)
	apiBoot := api.NewBoot(log, httpServer)
	if err := runner.Try(&apiBoot, 3, time.Second); err != nil {
		log.Fatal(err)
	}

	runner.WaitForTermination()
}
