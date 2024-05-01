package config

import (
	configError "accumulativeSystem/internal/errors/config"
	"errors"
	"flag"
	"os"
	"sync"
	"time"
)

// ConfigInstance - интерфейс для получения значений конфигурации
type ConfigInstance interface {
	GetRunAddress() string
	GetdatabaseURI() string
	GetAccrualSystemAddress() string
	GetMigrationPath() string
}

type Config struct {
	runAddress           string
	databaseURI          string
	accrualSystemAddress string
	migrationPath        string
	idleTimeout          time.Duration
}

var (
	instance ConfigInstance
	initOnce sync.Once
)

const (
	EnvLocal = "local"
	EnvTest  = "test"
	EnvProd  = "prod"
)

// MustLoad - функция для получения экземпляра ConfigInstance
func MustLoad() ConfigInstance {
	// initConfig является синглтоном, что для конфига не является критичным, так как он инициализируется один раз
	// и не будет больше меняться
	initOnce.Do(func() {
		var err error
		instance, err = initConfig()

		if err != nil {
			var configErr *configError.ConfigError

			if errors.As(err, &configErr) {
				panic(err)
			}

			panic(configError.New("failed to initialize config", err))
		}
	})

	return instance
}

func initConfig() (ConfigInstance, error) {
	var cfg Config

	if err := parseFlags(&cfg); err != nil {
		return nil, configError.New("error in parsing", err)
	}

	if err := overrideFromEnv(&cfg); err != nil {
		return nil, configError.New("error when overwriting env variables", err)
	}

	return &cfg, nil
}
func parseFlags(cfg *Config) error {
	//flag.StringVar(&cfg.runAddress, "a", "0.0.0.0:8080", "address and port to run server")
	//flag.StringVar(&cfg.databaseURI, "d", fmt.Sprintf("postgresql://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "postgres", "localhost", 5432, "postgres"), "")
	//flag.StringVar(&cfg.accrualSystemAddress, "r", "http://localhost:8081", "address of the accrual calculation system")

	flag.StringVar(&cfg.runAddress, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.databaseURI, "d", "", "")
	flag.StringVar(&cfg.accrualSystemAddress, "r", "", "address of the accrual calculation system")
	flag.StringVar(&cfg.migrationPath, "m", "file:./migrations", "")

	flag.Parse()
	return nil
}

func overrideFromEnv(cfg *Config) error {
	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		cfg.runAddress = envRunAddr
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		cfg.databaseURI = envDatabaseURI
	}

	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		cfg.accrualSystemAddress = envAccrualSystemAddress
	}

	return nil
}

// Реализация интерфейса ConfigInstance для структуры Config
func (c *Config) GetRunAddress() string {
	return c.runAddress
}

func (c *Config) GetdatabaseURI() string {
	return c.databaseURI
}

func (c *Config) GetAccrualSystemAddress() string {
	return c.accrualSystemAddress
}

func (c *Config) GetMigrationPath() string {
	return c.migrationPath
}
