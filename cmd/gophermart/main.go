package main

import (
	"accumulativeSystem/internal/config"
	"accumulativeSystem/internal/migrator"
)

func main() {
	config.LoadConfig()
	migrator.InitMigrator()
}
