package repository

import (
	"fmt"

	applog "github.com/bllooop/nameservice/internal/log"
	"github.com/jmoiron/sqlx"
)

type Config struct {
	Host       string
	Port       string
	Username   string
	Password   string
	DBname     string
	SSLMode    string
	ServerPort string
}

const (
	peopleListTable = "people"
)

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	applog.Logger.Info().Msg("Подключение к базе данных")
	constring := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.SSLMode)
	db, err := sqlx.Open("pgx", constring)
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBname, cfg.SSLMode)
	applog.Logger.Debug().Str("conn", connStr).Msg("Обработка подключения к БД")
	if err != nil {
		return nil, err
	}
	return db, nil
}
