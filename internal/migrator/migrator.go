package migrator

import (
	"accumulativeSystem/internal/config"
	migratorError "accumulativeSystem/internal/errors/migrator"
	_ "accumulativeSystem/internal/migrator/drivers" //инициализация нужных драйверов
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"sync"
)

var once sync.Once

func InitMigrator() {
	once.Do(func() {
		newMigrator()
	})
}

func newMigrator() {
	cnf := config.MustLoad()
	sourceUrl := cnf.MigrationPath
	databaseUrl := cnf.DatabaseUri

	instanceMigrate, err := migrate.New(sourceUrl, databaseUrl)

	if err != nil {
		//возвращаем панику, в случаи проблемы с инициализацией миграции
		panic(migratorError.New("error during migration", err))
	}

	if err := instanceMigrate.Up(); err != nil {
		if errors.As(err, &migrate.ErrNoChange) {
			return
		}

		//возвращаем панику, в случаи проблемы с применением миграции
		panic(migratorError.New("error during up migration", err))
	}
}