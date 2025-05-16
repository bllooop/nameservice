package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	logger.Log.Debug().Msg("Инициализация сервера...")

	if err := initConfig(); err != nil {
		logger.Log.Error().Err(err).Msg("")
		logger.Log.Fatal().Msg("Возникла ошибка загрузки конфига")
	}
	if err := godotenv.Load(); err != nil {
		logger.Log.Error().Err(err).Msg("")
		logger.Log.Fatal().Msg("Возникла ошибка с env")
	}
	logger.Log.Debug().Msg("Переменные окружения успешно загружены")
}

func initConfig() error {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
