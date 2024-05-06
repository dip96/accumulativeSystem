package app

import (
	"accumulativeSystem/internal/config"
	handBalance "accumulativeSystem/internal/http-server/handlers/balance/get"
	handWithdraw "accumulativeSystem/internal/http-server/handlers/balance/withdraw"
	handCreateOrder "accumulativeSystem/internal/http-server/handlers/order/create"
	handGetOrder "accumulativeSystem/internal/http-server/handlers/order/get"
	handLogin "accumulativeSystem/internal/http-server/handlers/user/login"
	handRegistration "accumulativeSystem/internal/http-server/handlers/user/registration"
	handUserWithdrawals "accumulativeSystem/internal/http-server/handlers/user/withdrawals"
	authMid "accumulativeSystem/internal/http-server/middleware/auth"
	"accumulativeSystem/internal/logger"
	"accumulativeSystem/internal/migrator"
	balanceRepository "accumulativeSystem/internal/repositories/balance"
	orderRepository "accumulativeSystem/internal/repositories/order"
	userRepository "accumulativeSystem/internal/repositories/user"
	balanceService "accumulativeSystem/internal/services/balance"
	orderService "accumulativeSystem/internal/services/order"
	orderQueue "accumulativeSystem/internal/services/order/queue"
	"accumulativeSystem/internal/services/order/queueService"
	userService "accumulativeSystem/internal/services/user"
	"accumulativeSystem/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"net/http"
	"time"
)

type App struct {
	Config   config.ConfigInstance
	Storage  storage.Storage
	Logger   logger.Logger
	Migrator migrator.Migrator
	JWTAuth  *jwtauth.JWTAuth
	Router   *chi.Mux
}

func NewApp(cfg config.ConfigInstance, storage storage.Storage, log logger.Logger, mig migrator.Migrator, jwtAuth *jwtauth.JWTAuth) (*App, error) {
	app := &App{
		Config:   cfg,
		Storage:  storage,
		Logger:   log,
		Migrator: mig,
		JWTAuth:  jwtAuth,
	}

	//run migrator
	err := mig.Up()

	if err != nil {
		app.Logger.Error("Error migrator")
		panic("Error migrator")
	}

	//Создаём каналы
	orderChan := orderQueue.NewOrderQueue()

	// Создаём репозиторий пользователей
	userRepo := userRepository.NewUserRepository(storage)
	orderRepo := orderRepository.NewOrderRepository(storage)
	balanceRepo := balanceRepository.NewBalanceRepository(storage)

	// Создаём сервис пользователей
	usService := userService.NewUserService(userRepo, balanceRepo, log)
	orService := orderService.NewOrderService(orderRepo, orderChan, log)
	baService := balanceService.NewBalanceService(balanceRepo, orderRepo, log)

	// Создаём каналы
	orderChanService := queueservice.NewOrderQueueService(cfg, orderChan, baService)

	//Run goroutine
	orderChanService.RunGoroutine(orService)

	//TODO добавить логер в роуты, сервисы, репозитории

	// Инициализация роута и хэндлеров
	//TODO добавить отдельный контейнер для хранения всех сервисов
	app.initRouter(usService, orService, baService)

	return app, nil
}

func (app *App) initRouter(userService userService.UserService, orderService orderService.OrderService, balanceService balanceService.BalanceService) {
	app.Router = chi.NewRouter()
	app.Router.Use(middleware.RequestID)
	app.Router.Use(middleware.Timeout(60 * time.Second))

	// Настройка публичных маршрутов
	app.Router.Group(func(r chi.Router) {
		r.Post("/api/user/register", handRegistration.New(userService, app.JWTAuth))
		r.Post("/api/user/login", handLogin.New(userService, app.JWTAuth))
	})

	// Настройка защищенных маршрутов
	app.Router.Group(func(r chi.Router) {
		//r.Use(authMid.JWTVerifier(app.JWTAuth))
		r.Use(authMid.AuthMiddleware(app.Logger))

		r.Post("/api/user/orders", handCreateOrder.New(orderService))
		r.Get("/api/user/orders", handGetOrder.New(orderService))
		r.Get("/api/user/balance", handBalance.New(balanceService))
		r.Post("/api/user/balance/withdraw", handWithdraw.New(balanceService))
		r.Get("/api/user/withdrawals", handUserWithdrawals.New(orderService))
	})
}

func (app *App) Run() error {
	srv := &http.Server{
		Addr:    app.Config.GetRunAddress(),
		Handler: app.Router,
	}

	// Обработка сигналов
	//signalChan := make(chan os.Signal, 1)
	////SIGINT - при попытке прервать процесс пользователем
	////SIGTERM - при попытке корректно прервать процесс средствами ос или системными утилитами
	//signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// горутина для ожидания сигналов SIGINT и SIGTERM
	//go func() {
	//	sig := <-signalChan
	//	app.Logger.Infof("Received signal: %s. Shutting down", sig)
	//
	//	app.Storage.Close()
	//	os.Exit(0)
	//}()

	app.Logger.Info("Starting server...")
	return srv.ListenAndServe()
}
