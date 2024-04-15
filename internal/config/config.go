package config

import (
	configError "accumulativeSystem/internal/errors/config"
	"errors"
	"flag"
	"fmt"
	"os"
	"sync"
	"time"
)

var (
	instanceConfig *Config
	initOnce       sync.Once
)

type Config struct {
	RunAddress           string
	DatabaseUri          string
	AccrualSystemAddress string
	MigrationPath        string
	IdleTimeout          time.Duration
}

type HttpServer struct {
	//todo разбить конфиги на типы
}

type Storage struct {
	//todo разбить конфиги на типы
}

func MustLoad() *Config {
	// initConfig является синглтоном, что для конфига не является критичным, так как он инициализируется один раз
	// и не будет больше меняться
	initOnce.Do(func() {
		var err error
		instanceConfig, err = initConfig()

		// Если возникла ошибка при инициализации конфигов, то приложение не будет работать корректно
		if err != nil {
			var configErr *configError.ConfigError

			if errors.As(err, &configErr) {
				panic(err)
			}

			// в случаи, если получили ошибку, при инициализации конфига, но она не нашего кастомного типа
			panic(configError.New("failed to initialize config", err))
		}
	})

	return instanceConfig
}

func initConfig() (*Config, error) {
	var cfg = Config{}

	if err := parseFlags(&cfg); err != nil {
		return nil, configError.New("error in parsing", err)
	}

	if err := overrideFromEnv(&cfg); err != nil {
		return nil, configError.New("error when overwriting env variables", err)
	}

	return &cfg, nil
}

func parseFlags(cfg *Config) error {
	//flag.StringVar(&cfg.RunAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.RunAddress, "a", "0.0.0.0:8080", "address and port to run server")
	flag.StringVar(&cfg.DatabaseUri, "d", fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "postgres", "localhost", 5432, "postgres"), "")
	flag.StringVar(&cfg.AccrualSystemAddress, "r", "", "File to save metrics")
	flag.StringVar(&cfg.MigrationPath, "m", "file:./migrations", "")

	flag.Parse()
	return nil
}

func overrideFromEnv(cfg *Config) error {
	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		cfg.RunAddress = envRunAddr
	}

	if envDatabaseUri := os.Getenv("DATABASE_URI"); envDatabaseUri != "" {
		cfg.DatabaseUri = envDatabaseUri
	}

	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		cfg.AccrualSystemAddress = envAccrualSystemAddress
	}

	return nil
}
