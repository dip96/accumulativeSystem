package migrator

import (
	"accumulativeSystem/internal/config"
	migratorError "accumulativeSystem/internal/errors/migrator"
	_ "accumulativeSystem/internal/migrator/drivers" //инициализация нужных драйверов
	"errors"
	"github.com/golang-migrate/migrate/v4"
)

type Migrator interface {
	Up() error
}

type migrator struct {
	instance *migrate.Migrate
}

func NewMigrator(cnf config.ConfigInstance) (Migrator, error) {
	sourceUrl := cnf.GetMigrationPath()
	databaseUrl := cnf.GetDatabaseURI()

	instanceMigrate, err := migrate.New(sourceUrl, databaseUrl)
	if err != nil {
		return nil, migratorError.New("error during migration", err)
	}

	return &migrator{instance: instanceMigrate}, nil
}

func (m *migrator) Up() error {
	if err := m.instance.Up(); err != nil {
		if errors.As(err, &migrate.ErrNoChange) {
			// изменений нет, можно выходить
			return nil
		}
		return migratorError.New("error during up migration", err)
	}

	return nil
}
