package main

import (
	"os"

	applog "github.com/bllooop/nameservice/internal/log"
	"github.com/bllooop/nameservice/internal/repository"
	running "github.com/bllooop/nameservice/internal/server"
	_ "github.com/jackc/pgx/v5/pgxpool"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	applog.Logger.Debug().Msg("Инициализация сервера...")

	if err := initConfig(); err != nil {
		applog.Logger.Error().Err(err).Msg("")
		applog.Logger.Fatal().Msg("Возникла ошибка загрузки конфига")
	}
	if err := godotenv.Load(); err != nil {
		applog.Logger.Error().Err(err).Msg("")
		applog.Logger.Fatal().Msg("Возникла ошибка с env")
	}
	applog.Logger.Debug().Msg("Переменные окружения успешно загружены")
	cfg := repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBname:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	running.Run(cfg)
}

func initConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
