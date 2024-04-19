package main

import (
	"accumulativeSystem/internal/config"
	handBalance "accumulativeSystem/internal/http-server/handlers/balance/getBalance"
	handCreateOrder "accumulativeSystem/internal/http-server/handlers/order/createOrder"
	handGetOrder "accumulativeSystem/internal/http-server/handlers/order/getOrders"
	handLogin "accumulativeSystem/internal/http-server/handlers/user/login"
	handRegistration "accumulativeSystem/internal/http-server/handlers/user/registration"
	authMid "accumulativeSystem/internal/http-server/middleware/auth"
	"accumulativeSystem/internal/logger"
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
	//TODO переделать на yaml?
	cnf := config.MustLoad()

	//TODO init logger
	//TODO add env local
	logger.Init("local")

	//TODO init migrator
	logger.Log.Info("init migrator")
	migrator.InitMigrator()

	//TODO init storage
	logger.Log.Info("init storage")
	storage := postgres.NewDb()

	//TODO jwt init
	logger.Log.Info("init jwt")
	jwtAuth := jwtauth.New("HS512", []byte("secret"), nil)

	//TODO router
	logger.Log.Info("init router")
	r := chi.NewRouter()
	r.Use(middleware.RequestID) //TODO добавить RequestID в Log.WithField()
	r.Use(middleware.Timeout(60 * time.Second))

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(jwtAuth))
		r.Use(authMid.AuthMiddleware)

		r.Post("/api/user/orders", handCreateOrder.New(storage))
		r.Get("/api/user/orders", handGetOrder.New(storage))
		r.Get("/api/user/balance", handBalance.New(storage))
	})

	// Public routes
	r.Group(func(r chi.Router) {
		r.Post("/api/user/register", handRegistration.New(storage, jwtAuth))
		r.Post("/api/user/login", handLogin.New(storage, jwtAuth))
	})

	srv := &http.Server{
		Addr:    cnf.RunAddress,
		Handler: r,
		//ReadTimeout:  cnf.Timeout,
		//WriteTimeout: cnf.Timeout,
		//IdleTimeout:  cnf.IdleTimeout,
	}

	logger.Log.Info("start server")
	if err := srv.ListenAndServe(); err != nil {
		//log.Error("failed to start server")
	}
}
