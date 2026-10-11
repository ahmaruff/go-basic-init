package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/ahmaruff/go-basic-init/internal/platform"
	"github.com/ahmaruff/go-basic-init/internal/platform/db"
	"github.com/ahmaruff/go-basic-init/internal/platform/logger"
	"github.com/joho/godotenv"
)

// version is stamped into every log line. Override at build time with
// -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		// The logger is not configured yet, so this reaches stdout only.
		logger.L().Error().Stack().Err(err).Msg("load config failed")
		return err
	}

	// start logger
	if logErr := setupLogger(cfg); logErr != nil {
		logger.L().Warn().Err(logErr).Msg("file logging disabled, stdout only")
	}

	defer logger.Close()

	database, err := openDatabase(cfg)
	if err != nil {
		logger.L().Error().Stack().Err(err).Msg("open database failed")
		return err
	}

	defer database.Close()

	return nil
}

func loadConfig() (platform.Config, error) {
	if err := godotenv.Load(); err != nil {
		return platform.Config{}, fmt.Errorf("load config: %w", err)
	}

	return platform.LoadConfig(), nil

}

func setupLogger(cfg platform.Config) error {

	return logger.Init(logger.Config{
		Env:         cfg.AppEnv,
		Level:       cfg.LogLevel,
		Service:     cfg.AppName,
		Version:     version,
		File:        cfg.LogFile,
		MaxSizeMB:   100,
		MaxBackups:  5,
		MaxAgeDays:  30,
		Compression: "gzip",
	})

}

func openDatabase(cfg platform.Config) (*sql.DB, error) {
	database, err := db.New(db.DBConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		Database: cfg.DBDatabase,
		Username: cfg.DBUsername,
		Password: cfg.DBPassword,
	})

	if err != nil {
		// The stack is captured where the error was born (db.New); here we only
		// add context, so formatting with %w keeps it intact.
		return nil, fmt.Errorf("connect database: %w", err)
	}

	return database, nil
}
