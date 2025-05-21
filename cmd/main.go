package main

import (
	"fmt"
	"os"

	applog "github.com/bllooop/nameservice/internal/log"
	"github.com/bllooop/nameservice/internal/repository"
	running "github.com/bllooop/nameservice/internal/server"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"github.com/spf13/viper"
)

// @title Name API
// @version 1.0
// @description API сервис по работе с фио и другими характеристиками

// @host localhost:8080
// @BasePath /

func main() {
	applog.InitLogger(os.Stdout, zerolog.DebugLevel)
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

	err := running.Run(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}

func initConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
