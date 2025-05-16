package repository

import (
	"database/sql"
	"fmt"

	applog "github.com/bllooop/nameservice/internal/log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose"
)

func RunMigrate(cfg Config, migratePath string) error {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.SSLMode)
	applog.Logger.Debug().Str("conn", connStr).Msg("Обработка подключения к БД")

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	defer db.Close()
	applog.Logger.Info().Msg("Применение миграций")
	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}
	err = goose.Up(db, migratePath)
	if err != nil {
		return err
	}
	applog.Logger.Info().Msg("Миграция прошла успешно!")
	return nil
}
