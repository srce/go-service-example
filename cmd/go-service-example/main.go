package main

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dzyanis/go-service-example/internal/api"
	"github.com/dzyanis/go-service-example/internal/config"
	"github.com/dzyanis/go-service-example/internal/migrations"
	transControllers "github.com/dzyanis/go-service-example/internal/transactions/controllers"
	transRepos "github.com/dzyanis/go-service-example/internal/transactions/repositories"
	transServices "github.com/dzyanis/go-service-example/internal/transactions/services"
	"github.com/dzyanis/go-service-example/internal/transactions/uow"
	usersControllers "github.com/dzyanis/go-service-example/internal/users/controllers"
	usersRepos "github.com/dzyanis/go-service-example/internal/users/repositories"
	usersServices "github.com/dzyanis/go-service-example/internal/users/services"
	walletsControllers "github.com/dzyanis/go-service-example/internal/wallets/controllers"
	walletsRepos "github.com/dzyanis/go-service-example/internal/wallets/repositories"
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
		log.Fatalf("loading config env: %v", err)
	}

	db := database.NewPostgres(cnf.Postgres)
	dbBoot := database.NewBoot(db)
	if err := runner.Try(&dbBoot, 10, time.Second); err != nil {
		log.Fatalf("new database instance: %v", err)
	}

	mg, err := migrations.NewPostgres(db.Connection())
	if err != nil {
		log.Fatalf("new migrations instance: %v", err)
	}
	migratBoot := migrations.NewBoot(log, mg)
	if err := runner.Once(&migratBoot); err != nil {
		log.Fatalf("booting migrations: %v", err)
	}

	jsonHelper := controllers.JSONHelper{}

	usersRepository := usersRepos.NewRepository(db)
	usersController := usersControllers.NewController(log,
		usersServices.NewService(usersRepository), jsonHelper)

	walletsRepository := walletsRepos.NewRepository(db)
	walletsController := walletsControllers.NewController(log,
		walletsServices.NewService(walletsRepository), jsonHelper)

	startUOW := func() (uow.UnitOfWork, error) {
		return uow.NewRepository(db.Connection())
	}

	transactionsRepository := transRepos.NewRepository(db)
	transactionsService := transServices.NewService(log,
		transactionsRepository, usersRepository, walletsRepository, startUOW)
	transactionsController := transControllers.NewController(log,
		transactionsService, jsonHelper)

	httpServer := api.NewServer(cnf.API,
		usersController, walletsController, transactionsController)
	apiBoot := api.NewBoot(log, httpServer)
	if err := runner.Try(&apiBoot, 3, time.Second); err != nil {
		log.Fatalf("booting http server: %v", err)
	}

	runner.WaitForTermination()
}
