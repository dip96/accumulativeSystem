package main

import (
	"accumulativeSystem/internal/config"
	handRegistration "accumulativeSystem/internal/http-server/handlers/user/registration"
	"accumulativeSystem/internal/migrator"
	"accumulativeSystem/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

func main() {
	//TODO init config
	cnf := config.MustLoad()

	//TODO init migrator
	migrator.InitMigrator()

	//TODO init logger

	//TODO init storage
	storage := postgres.NewDb()

	//TODO router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second))
	//chi.Use(logger.New())

	//r.Group()

	r.Post("/api/user/register", handRegistration.New(storage))

	srv := &http.Server{
		Addr:    cnf.RunAddress,
		Handler: r,
		//ReadTimeout:  cnf.Timeout,
		//WriteTimeout: cnf.Timeout,
		//IdleTimeout:  cnf.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		//log.Error("failed to start server")
	}
}
