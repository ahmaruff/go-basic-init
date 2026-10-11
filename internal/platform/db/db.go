package db

import (
	"database/sql"
	"fmt"

	"github.com/ahmaruff/go-basic-init/internal/platform/logger"
	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Host     string
	Port     string
	Database string
	Username string
	Password string
}

func New(cfg DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, logger.WithStack(fmt.Errorf("open mysql: %w", err))
	}

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, logger.WithStack(fmt.Errorf("ping mysql: %w", err))
	}

	return database, nil
}
