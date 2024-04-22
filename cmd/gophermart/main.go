package main

import (
	appInstance "accumulativeSystem/internal/app"
	"accumulativeSystem/internal/config"
	"accumulativeSystem/internal/logger"
	"accumulativeSystem/internal/migrator"
	"accumulativeSystem/internal/storage/postgres"
	"github.com/go-chi/jwtauth"
)

func main() {
	//init config
	//TODO переделать на yaml?
	cnf := config.MustLoad()

	//init logger
	log := logger.Init(config.EnvLocal)

	//init migrator
	log.Info("init migrator")
	migrator, err := migrator.NewMigrator(cnf)

	if err != nil {
		log.Error("failed to create migrator:", err)
		return
	}

	//init storage
	log.Info("init storage")
	storage, err := postgres.NewDb(cnf)

	if err != nil {
		log.Error("failed to create storage:", err)
		return
	}

	//jwt init
	log.Info("init jwt")
	jwtAuth := jwtauth.New("HS512", []byte("secret"), nil)

	app, err := appInstance.NewApp(cnf, storage, log, migrator, jwtAuth)

	if err != nil {
		log.Error("failed to create app instance:", err)
		return
	}

	if err := app.Run(); err != nil {
		log.Error("failed to run app: %v", err)
	}
}
