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
)

// @title Name API
// @version 1.0
// @description API сервис по работе с фио и другими характеристиками

// @host localhost:8080
// @BasePath /

func main() {
	applog.InitLogger(os.Stdout, zerolog.DebugLevel)
	applog.Logger.Debug().Msg("Инициализация сервера...")
	if err := godotenv.Load(); err != nil {
		applog.Logger.Error().Err(err).Msg("")
		applog.Logger.Fatal().Msg("Возникла ошибка с env")
	}
	applog.Logger.Debug().Msg("Переменные окружения успешно загружены")
	cfg := repository.Config{
		Host:       os.Getenv("HOST"),
		Port:       os.Getenv("PORT"),
		Username:   os.Getenv("USERNAME"),
		Password:   os.Getenv("DB_PASSWORD"),
		DBname:     os.Getenv("DBNAME"),
		SSLMode:    os.Getenv("SSLMODE"),
		ServerPort: os.Getenv("SERVERPORT"),
	}

	err := running.Run(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}
