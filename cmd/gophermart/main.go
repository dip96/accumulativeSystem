package main

import (
	"accumulativeSystem/internal/config"
	handLogin "accumulativeSystem/internal/http-server/handlers/user/login"
	handRegistration "accumulativeSystem/internal/http-server/handlers/user/registration"
	"accumulativeSystem/internal/migrator"
	"accumulativeSystem/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
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

	//TODO jwt init
	jwtAuth := jwtauth.New("HS512", []byte("secret"), nil)

	//TODO router
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Timeout(60 * time.Second))
	//chi.Use(logger.New())

	//r.Group()

	r.Post("/api/user/register", handRegistration.New(storage, jwtAuth))
	r.Post("/api/user/login", handLogin.New(storage, jwtAuth))

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
