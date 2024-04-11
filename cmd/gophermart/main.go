package main

import (
	"accumulativeSystem/internal/config"
	"accumulativeSystem/internal/migrator"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"time"
)

func main() {
	config.MustLoad()
	migrator.InitMigrator()

	//TODO router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second))
	//chi.Use(logger.New())

	//TODO init logger
}
